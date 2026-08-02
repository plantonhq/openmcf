package module

// webAuthBackendPy is the module-owned login backend delivered next to
// the user's locustfile when the web-UI login is on.
//
// WHY CODE: Locust's `--web-login` flag protects every web and REST
// route behind a session, but deliberately leaves the credential
// backend to locustfile code — the documented extension seam
// (extending-locust.html#authentication at the pin; the shape below
// follows upstream's own examples/web_ui_auth/basic.py). The chart's
// `master.auth` username/password values feed ONLY the legacy
// pre-2.21 code path that renders credentials as pod arguments —
// never engaged by this module.
//
// DELIVERY: this file rides the module-owned `<name>-web-auth`
// ConfigMap (mounted through the chart's extraConfigMaps seam) and is
// appended to the master's `-f` argument — the master loads it
// alongside the locustfile; workers never do (their default
// LOCUST_LOCUSTFILE env names the locustfile alone). Credentials and
// the Flask session-signing key arrive as Secret-projected FILES
// (mount_external_secret), never as environment or rendered values.
//
// PARITY: the Terraform module carries this exact content in its
// locals — keep both engines byte-identical.
const webAuthBackendPy = `"""Platform-managed login for the Locust web UI.

Locust protects its web routes behind a session when started with
--web-login, and delegates the credential backend to locustfile code.
This module implements a single-credential backend: the username,
password and the Flask session-signing key are read from
Secret-projected files, so nothing secret rides rendered values or
process arguments, and sessions survive pod restarts.
"""

import os
import secrets

from flask import Blueprint, redirect, request, session, url_for
from flask_login import UserMixin, login_user

from locust import events

_AUTH_DIR = "/opt/planton/web-auth"


def _read(name):
    with open(os.path.join(_AUTH_DIR, name)) as f:
        return f.read().strip()


class _AuthUser(UserMixin):
    def __init__(self, username):
        self.username = username

    def get_id(self):
        return self.username


@events.init.add_listener
def _planton_web_login(environment, **_kwargs):
    if not environment.web_ui:
        return

    username = _read("username")
    password = _read("password")

    web_ui = environment.web_ui
    web_ui.app.config["SECRET_KEY"] = _read("flask-secret-key")
    web_ui.login_manager.user_loader(
        lambda user_id: _AuthUser(user_id) if user_id == username else None
    )

    base_path = environment.parsed_options.web_base_path
    web_ui.auth_args = {
        "username_password_callback": base_path + "/planton/login",
    }

    blueprint = Blueprint("planton_auth", __name__, url_prefix=base_path)

    @blueprint.route("/planton/login", methods=["POST"])
    def _login():
        username_ok = secrets.compare_digest(
            request.form.get("username", ""), username
        )
        password_ok = secrets.compare_digest(
            request.form.get("password", ""), password
        )
        if username_ok and password_ok:
            login_user(_AuthUser(username))
            return redirect(url_for("locust.index"))
        session["auth_error"] = "Invalid username or password"
        return redirect(url_for("locust.login"))

    web_ui.app.register_blueprint(blueprint)
`
