'use client';

import React from 'react';
import { motion } from 'framer-motion';
import { Terminal, FileCode, GitBranch, Zap } from 'lucide-react';

export default function PlantonCLI() {
  const terminalCommands = `# Install planton CLI
brew install plantonhq/tap/planton

# Deploy using YAML
planton apply -f alb.yaml

# Check status
planton status

# View outputs
planton outputs`;

  const yamlExample = `apiVersion: aws.planton.dev/v1
kind: AwsAlb
metadata:
  env: dev
  name: dev-alb
  org: my-org
spec:
  idleTimeoutSeconds: 60
  securityGroups:
    - valueFrom:
        kind: AwsSecurityGroup
        name: dev-http-security-group
        fieldPath: status.outputs.securityGroupId
  ssl:
    certificateArn:
      value: arn:aws:acm:us-east-1:123456789:certificate/abc-123
    enabled: true
  subnets:
    - valueFrom:
        kind: AwsVpc
        name: dev-vpc
        fieldPath: status.outputs.publicSubnets.[0].id
    - valueFrom:
        kind: AwsVpc
        name: dev-vpc
        fieldPath: status.outputs.publicSubnets.[1].id`;

  return (
    <div className="h-full flex">
      {/* Left Side - Terminal (40%) */}
      <div className="w-[40%] bg-gradient-to-br from-gray-900 to-gray-800 p-8 overflow-y-auto border-r border-gray-700">
        <motion.div
          initial={{ opacity: 0, x: -20 }}
          animate={{ opacity: 1, x: 0 }}
          transition={{ duration: 0.6 }}
          className="space-y-6"
        >
          {/* Header */}
          <div>
            <div className="inline-flex items-center justify-center w-16 h-16 bg-gradient-to-br from-[#10b981] to-[#059669] rounded-2xl mb-4">
              <Terminal className="w-8 h-8 text-white" />
            </div>
            <h2 className="text-3xl font-black text-white mb-2">
              Deploy with the planton CLI
            </h2>
            <p className="text-gray-400">
              Simple CLI for YAML-based infrastructure deployment
            </p>
          </div>

          {/* Terminal Window */}
          <div className="bg-black rounded-2xl overflow-hidden border border-gray-700 demo-card-shadow">
            {/* Terminal Header */}
            <div className="bg-gray-800 px-4 py-2 flex items-center gap-2 border-b border-gray-700">
              <div className="flex gap-2">
                <div className="w-3 h-3 rounded-full bg-red-500"></div>
                <div className="w-3 h-3 rounded-full bg-[#333]"></div>
                <div className="w-3 h-3 rounded-full bg-[#10b981]"></div>
              </div>
              <div className="text-gray-400 text-xs ml-4">bash</div>
            </div>

            {/* Terminal Content */}
            <div className="p-6">
              <pre className="text-sm text-[#10b981] font-mono leading-loose">
                {terminalCommands}
              </pre>
            </div>
          </div>

          {/* CLI Features */}
          <div className="bg-gray-800/50 rounded-2xl p-6 border border-gray-700">
            <h3 className="font-bold text-white mb-4">CLI Features</h3>
            <div className="space-y-3">
              {[
                'Simple installation via Homebrew',
                'Apply infrastructure from YAML',
                'Check resource status',
                'View resource outputs',
                'Destroy resources cleanly',
              ].map((feature, i) => (
                <div key={i} className="flex items-start gap-3">
                  <div className="w-2 h-2 bg-[#10b981] rounded-full mt-2"></div>
                  <span className="text-sm text-gray-300">{feature}</span>
                </div>
              ))}
            </div>
          </div>
        </motion.div>
      </div>

      {/* Right Side - YAML Editor (60%) */}
      <div className="w-[60%] bg-gray-900 p-8 overflow-y-auto">
        <motion.div
          initial={{ opacity: 0, x: 20 }}
          animate={{ opacity: 1, x: 0 }}
          transition={{ duration: 0.6 }}
          className="space-y-6"
        >
          {/* Header */}
          <div>
            <div className="flex items-center justify-between mb-4">
              <h3 className="text-2xl font-bold text-white">alb.yaml</h3>
              <div className="px-3 py-1 bg-[#111] text-white text-xs font-bold rounded-full">
                YAML
              </div>
            </div>
            <p className="text-gray-400 mb-6">
              Declarative infrastructure specification
            </p>
          </div>

          {/* YAML Editor */}
          <div className="bg-gray-800 rounded-xl p-6 border border-gray-700">
            <pre className="text-sm text-gray-100 font-mono leading-relaxed">
              <code>{yamlExample}</code>
            </pre>
          </div>

          {/* Callouts */}
          <div className="space-y-4">
            <div className="bg-[#1a1a1a]/30 border border-[#3a3a3a] rounded-xl p-6">
              <h4 className="text-white font-bold mb-2 flex items-center gap-2">
                <GitBranch className="w-5 h-5" />
                GitOps-Ready, Declarative Infrastructure
              </h4>
              <p className="text-gray-300 text-sm leading-relaxed">
                Commit your YAML to Git, deploy via CI/CD. Full version control and audit trail for your infrastructure.
              </p>
            </div>

            <div className="bg-[#111]/30 border border-[#3a3a3a] rounded-xl p-6">
              <h4 className="text-white font-bold mb-2 flex items-center gap-2">
                <Zap className="w-5 h-5" />
                Reference Other Resources with valueFrom
              </h4>
              <p className="text-gray-300 text-sm leading-relaxed">
                Link resources together dynamically. Security groups, subnets, and other outputs can be referenced automatically.
              </p>
            </div>
          </div>

          {/* Features */}
          <div className="space-y-3">
            <h4 className="text-white font-bold">YAML Benefits</h4>
            {[
              { label: 'Type-safe validation', color: 'bg-[#333]' },
              { label: 'Git-friendly format', color: 'bg-[#10b981]' },
              { label: 'Resource references', color: 'bg-[#333]' },
              { label: 'Easy to read and review', color: 'bg-[#333]' },
              { label: 'CI/CD integration', color: 'bg-[#333]' },
            ].map((item, i) => (
              <div key={i} className="flex items-start gap-3">
                <div className={`w-2 h-2 ${item.color} rounded-full mt-2`}></div>
                <div>
                  <div className="text-white font-medium">{item.label}</div>
                </div>
              </div>
            ))}
          </div>

          {/* Bottom CTA */}
          <div className="bg-gradient-to-r from-[#1a1a1a]/50 to-[#1a1a1a]/50 border border-[#10b981] rounded-xl p-6">
            <h4 className="text-white font-bold mb-2 flex items-center gap-2">
              <FileCode className="w-5 h-5" />
              Open Source & Free
            </h4>
            <p className="text-gray-300 text-sm leading-relaxed">
              The CLI and all modules are completely open source. Use it for free, customize it, contribute back.
            </p>
          </div>
        </motion.div>
      </div>
    </div>
  );
}
