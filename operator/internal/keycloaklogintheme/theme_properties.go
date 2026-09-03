package keycloaklogintheme

// themeProperties is the theme's manifest file: Keycloak reads it to learn
// what this theme extends and which stylesheets to serve.
//
// parent=keycloak.v2 is the load-bearing decision. Extending the bundled
// theme means every login-family screen Keycloak can ever render -- OTP
// enrollment, recovery codes, WebAuthn, future LDAP-federation screens --
// keeps working with the parent's templates and PatternFly styling, and this
// theme only lays Planton's palette on top. A from-scratch theme
// (parent=base) would leave any screen we forgot as unstyled raw HTML.
//
// styles: entries resolve against THIS theme first, then the parent -- so
// css/styles.css serves the parent's own stylesheet (small: logo/layout
// helpers) and css/planton.css adds ours after it, winning ties by order.
// stylesCommon (the PatternFly v5 vendor stylesheets) is inherited from the
// parent and deliberately not redeclared.
//
// darkMode=false pins the palette. The parent ships an OS-preference
// listener that toggles PatternFly's dark class; Planton's login is dark BY
// DESIGN (matching the console's default), so the OS must not be able to
// flip it to a light PatternFly look this theme never styled.
const themeProperties = `parent=keycloak.v2
import=common/keycloak

styles=css/styles.css css/planton.css

darkMode=false
`
