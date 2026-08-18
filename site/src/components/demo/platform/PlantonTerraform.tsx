'use client';

import React from 'react';
import { motion } from 'framer-motion';
import { FolderTree, FileCode, CheckCircle2 } from 'lucide-react';

export default function PlantonTerraform() {
  const terraformCode = `resource "aws_lb" "main" {
  name               = var.alb_name
  internal           = var.internal
  load_balancer_type = "application"
  security_groups    = var.security_groups
  subnets            = var.subnets
  
  idle_timeout = var.idle_timeout_seconds
  
  enable_deletion_protection = var.enable_deletion_protection
  enable_http2              = var.enable_http2
  
  tags = merge(
    var.tags,
    {
      "ManagedBy" = "planton"
    }
  )
}`;

  return (
    <div className="h-full flex">
      {/* Left Side - Directory Structure (40%) */}
      <div className="w-[40%] bg-gradient-to-br from-[#151515] to-[#151515] p-8 overflow-y-auto border-r border-[#2a2a2a]">
        <motion.div
          initial={{ opacity: 0, x: -20 }}
          animate={{ opacity: 1, x: 0 }}
          transition={{ duration: 0.6 }}
          className="space-y-6"
        >
          {/* Header */}
          <div>
            <div className="inline-flex items-center justify-center w-16 h-16 bg-gradient-to-br from-white to-[#666] rounded-2xl mb-4">
              <FolderTree className="w-8 h-8 text-white" />
            </div>
            <h2 className="text-3xl font-black text-gray-900 mb-2">
              Production-Ready Terraform Modules
            </h2>
            <p className="text-gray-600">
              Pre-written, battle-tested infrastructure code
            </p>
          </div>

          {/* Directory Tree */}
          <div className="bg-white rounded-2xl p-6 demo-card-shadow border-2 border-[#2a2a2a]">
            <h3 className="font-bold text-gray-900 mb-4 flex items-center gap-2">
              <FileCode className="w-5 h-5" />
              Repository Structure
            </h3>
            <div className="font-mono text-sm text-gray-700 space-y-1">
              <div>planton/</div>
              <div className="ml-4">├── apis/</div>
              <div className="ml-8">│   └── aws/alb/v1/</div>
              <div className="ml-12 text-white font-semibold">│       └── spec.yaml</div>
              <div className="ml-4">└── modules/</div>
              <div className="ml-8">    └── aws-alb/</div>
              <div className="ml-12 text-white font-semibold">        ├── main.tf</div>
              <div className="ml-12 text-white font-semibold">        ├── variables.tf</div>
              <div className="ml-12 text-white font-semibold">        └── outputs.tf</div>
            </div>
          </div>

          {/* Features */}
          <div className="bg-white rounded-2xl p-6 demo-card-shadow space-y-4">
            <h3 className="font-bold text-gray-900 mb-4">What&apos;s Included</h3>
            {[
              'Production-grade Terraform',
              'Best practices built-in',
              'Comprehensive variable definitions',
              'Useful output values',
              'Proper error handling',
              'Security configurations',
            ].map((feature, i) => (
              <div key={i} className="flex items-start gap-3">
                <CheckCircle2 className="w-5 h-5 text-[#10b981] flex-shrink-0 mt-0.5" />
                <span className="text-sm text-gray-700">{feature}</span>
              </div>
            ))}
          </div>
        </motion.div>
      </div>

      {/* Right Side - Code Preview (60%) */}
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
              <h3 className="text-2xl font-bold text-white">main.tf</h3>
              <div className="px-3 py-1 bg-[#111] text-white text-xs font-bold rounded-full">
                aws-alb module
              </div>
            </div>
            <p className="text-gray-400 mb-6">
              Example Terraform module for AWS Application Load Balancer
            </p>
          </div>

          {/* Code Block */}
          <div className="bg-gray-800 rounded-xl p-6 border border-gray-700">
            <pre className="text-sm text-gray-100 font-mono leading-relaxed">
              <code>{terraformCode}</code>
            </pre>
          </div>

          {/* Callouts */}
          <div className="space-y-4">
            <div className="bg-[#111]/30 border border-[#3a3a3a] rounded-xl p-6">
              <h4 className="text-white font-bold mb-2">Don&apos;t Write Terraform from Scratch</h4>
              <p className="text-gray-300 text-sm leading-relaxed">
                Every component comes with production-ready Terraform modules. 
                Just provide your values—we&apos;ve handled the complexity.
              </p>
            </div>

            <div className="bg-[#1a1a1a]/30 border border-[#10b981] rounded-xl p-6">
              <h4 className="text-white font-bold mb-2">Battle-Tested, Production-Ready</h4>
              <p className="text-gray-300 text-sm leading-relaxed">
                These modules are used in production by real companies. 
                Security hardened, optimized, and continuously updated.
              </p>
            </div>
          </div>

          {/* Features List */}
          <div className="space-y-3">
            <h4 className="text-white font-bold">Module Features</h4>
            {[
              { label: 'Security best practices', color: 'bg-red-400' },
              { label: 'Sensible defaults', color: 'bg-[#333]' },
              { label: 'Customizable variables', color: 'bg-[#10b981]' },
              { label: 'Comprehensive outputs', color: 'bg-[#333]' },
              { label: 'Error handling', color: 'bg-[#333]' },
            ].map((item, i) => (
              <div key={i} className="flex items-start gap-3">
                <div className={`w-2 h-2 ${item.color} rounded-full mt-2`}></div>
                <div>
                  <div className="text-white font-medium">{item.label}</div>
                </div>
              </div>
            ))}
          </div>
        </motion.div>
      </div>
    </div>
  );
}
