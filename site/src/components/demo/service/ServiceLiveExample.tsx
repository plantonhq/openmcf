'use client';

import React from 'react';
import { motion } from 'framer-motion';
import { ExternalLink, CheckCircle, Zap, Clock, TrendingUp } from 'lucide-react';

export default function ServiceLiveExample() {
  return (
    <div className="h-full flex flex-col">
      {/* Header */}
      <div className="px-8 py-6 border-b border-gray-200">
        <h2 className="text-2xl font-bold text-gray-900">Live Service Example</h2>
        <p className="text-gray-600 mt-1">
          Your service is now live and serving production traffic
        </p>
      </div>

      <div className="flex-1 overflow-auto bg-gradient-to-br from-[#151515] via-[#151515] to-[#151515]">
        <div className="max-w-5xl mx-auto p-8">
          {/* Service card */}
          <motion.div
            initial={{ opacity: 0, y: 30 }}
            animate={{ opacity: 1, y: 0 }}
            className="bg-white rounded-2xl shadow-2xl overflow-hidden border border-gray-200 mb-8"
          >
            {/* Header */}
            <div className="bg-gradient-to-r from-white to-[#666] px-8 py-6 text-white">
              <div className="flex items-center justify-between mb-2">
                <h3 className="text-2xl font-bold">Billing Service</h3>
                <div className="flex items-center gap-2 bg-white/20 backdrop-blur-sm px-3 py-1 rounded-full">
                  <div className="w-2 h-2 bg-[#10b981] rounded-full animate-pulse" />
                  <span className="text-sm font-medium">Live</span>
                </div>
              </div>
              <p className="text-gray-300">Java Spring Boot - Payment Processing API</p>
            </div>

            {/* Content */}
            <div className="p-8">
              {/* Live URL */}
              <div className="mb-8">
                <label className="block text-sm font-semibold text-gray-700 mb-3">
                  Live URL
                </label>
                <div className="flex items-center gap-3">
                  <div className="flex-1 bg-gray-50 border border-gray-300 rounded-lg px-4 py-3 font-mono text-white text-sm font-semibold">
                    https://billing-api.acmecorp.com
                  </div>
                  <button className="px-4 py-3 bg-gradient-to-r from-white to-[#666] text-white rounded-lg font-medium hover:from-white hover:to-[#666] transition-colors flex items-center gap-2 shadow-lg">
                    <ExternalLink className="w-4 h-4" />
                    Visit
                  </button>
                </div>
              </div>

              {/* Deployment details */}
              <div className="grid md:grid-cols-2 gap-6 mb-8">
                <div className="bg-[#151515] border border-[#2a2a2a] rounded-xl p-6">
                  <div className="flex items-center gap-3 mb-3">
                    <CheckCircle className="w-5 h-5 text-[#10b981]" />
                    <h4 className="font-bold text-gray-900">No Dockerfile</h4>
                  </div>
                  <p className="text-sm text-gray-600">
                    Deployed using BuildPacks with automatic Java runtime detection and optimization.
                  </p>
                </div>

                <div className="bg-[#151515] border border-[#2a2a2a] rounded-xl p-6">
                  <div className="flex items-center gap-3 mb-3">
                    <CheckCircle className="w-5 h-5 text-white" />
                    <h4 className="font-bold text-gray-900">No GitHub Actions</h4>
                  </div>
                  <p className="text-sm text-gray-600">
                    Platform-managed CI/CD handles all builds and deployments automatically.
                  </p>
                </div>
              </div>

              {/* Metrics */}
              <div className="grid grid-cols-4 gap-4 mb-8">
                <motion.div
                  initial={{ opacity: 0, y: 20 }}
                  animate={{ opacity: 1, y: 0 }}
                  transition={{ delay: 0.1 }}
                  className="bg-gradient-to-br from-[#151515] to-[#1a1a1a] border border-[#2a2a2a] rounded-xl p-4"
                >
                  <div className="flex items-center gap-2 mb-2">
                    <TrendingUp className="w-4 h-4 text-white" />
                    <span className="text-xs font-semibold text-white">Requests</span>
                  </div>
                  <div className="text-2xl font-bold text-gray-900">1.2M</div>
                  <div className="text-xs text-gray-600">Last 24h</div>
                </motion.div>

                <motion.div
                  initial={{ opacity: 0, y: 20 }}
                  animate={{ opacity: 1, y: 0 }}
                  transition={{ delay: 0.2 }}
                  className="bg-gradient-to-br from-[#151515] to-[#1a1a1a] border border-[#2a2a2a] rounded-xl p-4"
                >
                  <div className="flex items-center gap-2 mb-2">
                    <Zap className="w-4 h-4 text-white" />
                    <span className="text-xs font-semibold text-white">Latency</span>
                  </div>
                  <div className="text-2xl font-bold text-gray-900">45ms</div>
                  <div className="text-xs text-gray-600">p95</div>
                </motion.div>

                <motion.div
                  initial={{ opacity: 0, y: 20 }}
                  animate={{ opacity: 1, y: 0 }}
                  transition={{ delay: 0.3 }}
                  className="bg-gradient-to-br from-[#151515] to-[#1a1a1a] border border-[#2a2a2a] rounded-xl p-4"
                >
                  <div className="flex items-center gap-2 mb-2">
                    <CheckCircle className="w-4 h-4 text-[#10b981]" />
                    <span className="text-xs font-semibold text-[#10b981]">Uptime</span>
                  </div>
                  <div className="text-2xl font-bold text-gray-900">99.9%</div>
                  <div className="text-xs text-gray-600">30 days</div>
                </motion.div>

                <motion.div
                  initial={{ opacity: 0, y: 20 }}
                  animate={{ opacity: 1, y: 0 }}
                  transition={{ delay: 0.4 }}
                  className="bg-gradient-to-br from-[#151515] to-[#1a1a1a] border border-[#2a2a2a] rounded-xl p-4"
                >
                  <div className="flex items-center gap-2 mb-2">
                    <Clock className="w-4 h-4 text-[#a0a0a0]" />
                    <span className="text-xs font-semibold text-[#a0a0a0]">Deploys</span>
                  </div>
                  <div className="text-2xl font-bold text-gray-900">47</div>
                  <div className="text-xs text-gray-600">This month</div>
                </motion.div>
              </div>

              {/* Configuration summary */}
              <div className="bg-gray-50 border border-gray-200 rounded-xl p-6">
                <h4 className="font-bold text-gray-900 mb-4">Service Configuration</h4>
                <div className="grid md:grid-cols-2 gap-6 text-sm">
                  <div>
                    <div className="text-gray-500 mb-1">Repository</div>
                    <div className="text-gray-900 font-mono">acmecorp-engineering/billing-service</div>
                  </div>
                  <div>
                    <div className="text-gray-500 mb-1">Branch</div>
                    <div className="text-gray-900 font-mono">main</div>
                  </div>
                  <div>
                    <div className="text-gray-500 mb-1">Build Method</div>
                    <div className="text-gray-900">BuildPacks (Java)</div>
                  </div>
                  <div>
                    <div className="text-gray-500 mb-1">Deployment Target</div>
                    <div className="text-gray-900">AWS EKS (us-east-1)</div>
                  </div>
                  <div>
                    <div className="text-gray-500 mb-1">Secrets Backend</div>
                    <div className="text-gray-900">AWS Secrets Manager</div>
                  </div>
                  <div>
                    <div className="text-gray-500 mb-1">SSL Certificate</div>
                    <div className="text-gray-900">Auto-managed</div>
                  </div>
                </div>
              </div>
            </div>
          </motion.div>

        </div>
      </div>
    </div>
  );
}

