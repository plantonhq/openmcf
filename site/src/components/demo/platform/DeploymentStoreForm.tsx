'use client';

import React, { useState } from 'react';
import { motion } from 'framer-motion';
import { Play, Code } from 'lucide-react';

export default function DeploymentStoreForm() {
  const [formData, setFormData] = useState({
    env: 'dev',
    name: 'dev-alb',
    subnet1: 'subnet-abc123',
    subnet2: 'subnet-def456',
    securityGroup: 'sg-xyz789',
    sslEnabled: true,
    idleTimeout: '60',
  });

  const yamlPreview = `apiVersion: aws.planton.dev/v1alpha1
kind: AwsAlb
metadata:
  env: ${formData.env}
  name: ${formData.name}
spec:
  idleTimeoutSeconds: ${formData.idleTimeout}
  ssl:
    enabled: ${formData.sslEnabled}
  subnets:
    - ${formData.subnet1}
    - ${formData.subnet2}
  securityGroups:
    - ${formData.securityGroup}`;

  return (
    <div className="h-full flex">
      {/* Left Side - Form (40%) */}
      <div className="w-[40%] bg-gradient-to-br from-[#151515] to-[#151515] p-8 overflow-y-auto border-r border-[#2a2a2a]">
        <motion.div
          initial={{ opacity: 0, x: -20 }}
          animate={{ opacity: 1, x: 0 }}
          transition={{ duration: 0.6 }}
          className="space-y-6"
        >
          {/* Header */}
          <div>
            <h2 className="text-3xl font-black text-gray-900 mb-2">
              Deploy AWS ALB
            </h2>
            <p className="text-gray-600">
              Configure via form - YAML generated automatically
            </p>
          </div>

          {/* Form Fields */}
          <div className="bg-white rounded-2xl p-6 demo-card-shadow space-y-4">
            <div>
              <label className="block text-sm font-semibold text-gray-700 mb-2">
                Environment
              </label>
              <input
                type="text"
                value={formData.env}
                onChange={(e) => setFormData({ ...formData, env: e.target.value })}
                className="w-full px-4 py-2 border-2 border-gray-200 rounded-lg focus:border-[#3a3a3a] focus:outline-none"
              />
            </div>

            <div>
              <label className="block text-sm font-semibold text-gray-700 mb-2">
                Name
              </label>
              <input
                type="text"
                value={formData.name}
                onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                className="w-full px-4 py-2 border-2 border-gray-200 rounded-lg focus:border-[#3a3a3a] focus:outline-none"
              />
            </div>

            <div>
              <label className="block text-sm font-semibold text-gray-700 mb-2">
                Subnet 1
              </label>
              <select
                value={formData.subnet1}
                onChange={(e) => setFormData({ ...formData, subnet1: e.target.value })}
                className="w-full px-4 py-2 border-2 border-gray-200 rounded-lg focus:border-[#3a3a3a] focus:outline-none"
              >
                <option value="subnet-abc123">subnet-abc123 (us-east-1a)</option>
                <option value="subnet-def456">subnet-def456 (us-east-1b)</option>
              </select>
            </div>

            <div>
              <label className="block text-sm font-semibold text-gray-700 mb-2">
                Subnet 2
              </label>
              <select
                value={formData.subnet2}
                onChange={(e) => setFormData({ ...formData, subnet2: e.target.value })}
                className="w-full px-4 py-2 border-2 border-gray-200 rounded-lg focus:border-[#3a3a3a] focus:outline-none"
              >
                <option value="subnet-def456">subnet-def456 (us-east-1b)</option>
                <option value="subnet-ghi789">subnet-ghi789 (us-east-1c)</option>
              </select>
            </div>

            <div>
              <label className="block text-sm font-semibold text-gray-700 mb-2">
                Security Group
              </label>
              <select
                value={formData.securityGroup}
                onChange={(e) => setFormData({ ...formData, securityGroup: e.target.value })}
                className="w-full px-4 py-2 border-2 border-gray-200 rounded-lg focus:border-[#3a3a3a] focus:outline-none"
              >
                <option value="sg-xyz789">sg-xyz789 (http-https)</option>
                <option value="sg-abc123">sg-abc123 (all-traffic)</option>
              </select>
            </div>

            <div className="flex items-center gap-3">
              <input
                type="checkbox"
                id="ssl"
                checked={formData.sslEnabled}
                onChange={(e) => setFormData({ ...formData, sslEnabled: e.target.checked })}
                className="w-5 h-5 text-[#a0a0a0] border-2 border-gray-300 rounded focus:ring-white"
              />
              <label htmlFor="ssl" className="text-sm font-semibold text-gray-700">
                Enable SSL/TLS
              </label>
            </div>

            <div>
              <label className="block text-sm font-semibold text-gray-700 mb-2">
                Idle Timeout (seconds)
              </label>
              <input
                type="number"
                value={formData.idleTimeout}
                onChange={(e) => setFormData({ ...formData, idleTimeout: e.target.value })}
                className="w-full px-4 py-2 border-2 border-gray-200 rounded-lg focus:border-[#3a3a3a] focus:outline-none"
              />
            </div>
          </div>

          {/* Deploy Button */}
          <button className="w-full bg-gradient-to-r from-[#999] to-red-600 text-white px-8 py-4 rounded-xl font-bold text-lg hover:from-[#888] hover:to-red-700 transition-all shadow-lg hover:shadow-xl flex items-center justify-center gap-2">
            <Play className="w-5 h-5" />
            Deploy ALB
          </button>

          {/* Alternative */}
          <div className="text-center">
            <button className="text-sm text-[#a0a0a0] hover:text-[#a0a0a0] font-medium flex items-center gap-2 mx-auto">
              <Code className="w-4 h-4" />
              Use YAML instead
            </button>
          </div>
        </motion.div>
      </div>

      {/* Right Side - YAML Preview (60%) */}
      <div className="w-[60%] bg-gray-900 p-8 overflow-y-auto">
        <motion.div
          initial={{ opacity: 0, x: 20 }}
          animate={{ opacity: 1, x: 0 }}
          transition={{ duration: 0.6 }}
          className="space-y-4"
        >
          {/* Header */}
          <div className="flex items-center justify-between mb-4">
            <h3 className="text-xl font-bold text-white">Live YAML Preview</h3>
            <div className="px-3 py-1 bg-[#10b981] text-white text-xs font-bold rounded-full">
              Auto-generated
            </div>
          </div>

          {/* YAML Code Block */}
          <div className="bg-gray-800 rounded-xl p-6 border border-gray-700">
            <pre className="text-sm text-gray-100 font-mono leading-relaxed">
              {yamlPreview}
            </pre>
          </div>

          {/* Info Box */}
          <div className="bg-[#1a1a1a]/30 border border-[#3a3a3a] rounded-xl p-6">
            <h4 className="text-white font-bold mb-2 flex items-center gap-2">
              <Code className="w-5 h-5" />
              Forms or YAML - Your Choice
            </h4>
            <p className="text-gray-300 text-sm leading-relaxed">
              Use the form for quick deployments, or edit the YAML directly for advanced configurations. 
              Both methods generate the same infrastructure.
            </p>
          </div>

          {/* Features */}
          <div className="space-y-3">
            <div className="flex items-start gap-3">
              <div className="w-2 h-2 bg-[#10b981] rounded-full mt-2"></div>
              <div>
                <div className="text-white font-medium">Type-safe YAML</div>
                <div className="text-gray-400 text-sm">Validated against component schemas</div>
              </div>
            </div>
            <div className="flex items-start gap-3">
              <div className="w-2 h-2 bg-[#333] rounded-full mt-2"></div>
              <div>
                <div className="text-white font-medium">GitOps ready</div>
                <div className="text-gray-400 text-sm">Commit to version control, deploy with CI/CD</div>
              </div>
            </div>
            <div className="flex items-start gap-3">
              <div className="w-2 h-2 bg-[#333] rounded-full mt-2"></div>
              <div>
                <div className="text-white font-medium">Resource references</div>
                <div className="text-gray-400 text-sm">Link to other resources with valueFrom</div>
              </div>
            </div>
          </div>
        </motion.div>
      </div>
    </div>
  );
}

