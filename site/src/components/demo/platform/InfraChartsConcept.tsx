'use client';

import React from 'react';
import { motion } from 'framer-motion';
import { Package, GitBranch, Zap, ChevronRight } from 'lucide-react';

export default function InfraChartsConcept() {
  return (
    <div className="h-full overflow-y-auto">
      <div className="max-w-7xl mx-auto p-8 space-y-8">
        {/* Header */}
        <motion.div
          initial={{ opacity: 0, y: -20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.6 }}
          className="text-center"
        >
          <h1 className="text-4xl md:text-5xl font-black text-gray-900 mb-4">
            Infra Charts
          </h1>
          <div className="inline-block bg-gradient-to-r from-white to-[#666] text-white px-6 py-3 rounded-full mb-4">
            <p className="text-lg font-bold">
              Helm Charts for Kubernetes. Infra Charts for Cloud.
            </p>
          </div>
          <p className="text-xl text-gray-600 max-w-3xl mx-auto">
            Package multiple Lego Blocks into complete, production-ready environments
          </p>
        </motion.div>

        {/* The Inspiration */}
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.2, duration: 0.6 }}
          className="bg-gradient-to-r from-[#151515] to-[#151515] border-2 border-[#2a2a2a] rounded-2xl p-8"
        >
          <h2 className="text-2xl font-bold text-gray-900 mb-4 flex items-center gap-3">
            <div className="w-10 h-10 bg-[#111] rounded-lg flex items-center justify-center">
              <Package className="w-6 h-6 text-white" />
            </div>
            The Inspiration: Helm Charts
          </h2>
          <div className="grid md:grid-cols-2 gap-6">
            <div className="bg-white rounded-xl p-6 border border-[#2a2a2a]">
              <h3 className="font-bold text-lg mb-3 text-white">Kubernetes Helm Charts</h3>
              <p className="text-gray-700 mb-3">
                Package multiple Kubernetes resources together as reusable templates
              </p>
              <div className="space-y-2 text-sm text-gray-600">
                <div className="flex items-center gap-2">
                  <ChevronRight className="w-4 h-4 text-white" />
                  <span>Deploy complete applications</span>
                </div>
                <div className="flex items-center gap-2">
                  <ChevronRight className="w-4 h-4 text-white" />
                  <span>Parameterized templates</span>
                </div>
                <div className="flex items-center gap-2">
                  <ChevronRight className="w-4 h-4 text-white" />
                  <span>Reusable across environments</span>
                </div>
              </div>
            </div>

            <div className="bg-gradient-to-br from-white to-[#666] rounded-xl p-6 text-white">
              <h3 className="font-bold text-lg mb-3">Infra Charts</h3>
              <p className="mb-3 text-white/90">
                Package multiple Lego Blocks together as reusable infrastructure templates
              </p>
              <div className="space-y-2 text-sm">
                <div className="flex items-center gap-2">
                  <ChevronRight className="w-4 h-4" />
                  <span>Deploy complete environments</span>
                </div>
                <div className="flex items-center gap-2">
                  <ChevronRight className="w-4 h-4" />
                  <span>Parameterized templates</span>
                </div>
                <div className="flex items-center gap-2">
                  <ChevronRight className="w-4 h-4" />
                  <span>Reusable across projects</span>
                </div>
              </div>
            </div>
          </div>
        </motion.div>

        {/* What It Does */}
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.4, duration: 0.6 }}
          className="bg-[#151515] border-2 border-[#2a2a2a] rounded-2xl p-8"
        >
          <h2 className="text-2xl font-bold text-gray-900 mb-6 flex items-center gap-3">
            <div className="w-10 h-10 bg-[#10b981] rounded-lg flex items-center justify-center">
              <Zap className="w-6 h-6 text-white" />
            </div>
            What Infra Charts Do
          </h2>
          <div className="grid md:grid-cols-3 gap-6">
            <div className="bg-white rounded-xl p-6 border border-[#2a2a2a]">
              <GitBranch className="w-8 h-8 text-[#10b981] mb-3" />
              <h3 className="font-bold mb-2">Orchestrated Deployment</h3>
              <p className="text-sm text-gray-700">
                Automatically handles dependencies and deploys resources in the correct order using a DAG
              </p>
            </div>
            <div className="bg-white rounded-xl p-6 border border-[#2a2a2a]">
              <Package className="w-8 h-8 text-[#10b981] mb-3" />
              <h3 className="font-bold mb-2">Complete Environments</h3>
              <p className="text-sm text-gray-700">
                Package all resources needed for production: VPC, security, compute, storage, DNS, SSL
              </p>
            </div>
            <div className="bg-white rounded-xl p-6 border border-[#2a2a2a]">
              <Zap className="w-8 h-8 text-[#10b981] mb-3" />
              <h3 className="font-bold mb-2">20 Minutes to Production</h3>
              <p className="text-sm text-gray-700">
                What used to take weeks of manual work now deploys in ~20 minutes
              </p>
            </div>
          </div>
        </motion.div>

        {/* How It Works */}
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.6, duration: 0.6 }}
          className="bg-gradient-to-r from-[#151515] to-[#151515] border-2 border-[#2a2a2a] rounded-2xl p-8"
        >
          <h2 className="text-2xl font-bold text-gray-900 mb-6 text-center">
            How It Works
          </h2>
          <div className="flex flex-col md:flex-row items-center justify-center gap-4">
            <div className="bg-white rounded-xl p-6 border-2 border-[#3a3a3a] text-center flex-1 max-w-xs">
              <div className="text-3xl font-black text-white mb-2">1</div>
              <div className="font-semibold mb-2">Fill Simple Form</div>
              <div className="text-sm text-gray-600">7 parameters</div>
            </div>
            <ChevronRight className="w-8 h-8 text-white rotate-90 md:rotate-0" />
            <div className="bg-white rounded-xl p-6 border-2 border-[#3a3a3a] text-center flex-1 max-w-xs">
              <div className="text-3xl font-black text-white mb-2">2</div>
              <div className="font-semibold mb-2">Template Rendering</div>
              <div className="text-sm text-gray-600">Values injected into chart</div>
            </div>
            <ChevronRight className="w-8 h-8 text-white rotate-90 md:rotate-0" />
            <div className="bg-white rounded-xl p-6 border-2 border-[#3a3a3a] text-center flex-1 max-w-xs">
              <div className="text-3xl font-black text-white mb-2">3</div>
              <div className="font-semibold mb-2">DAG Deployment</div>
              <div className="text-sm text-gray-600">9 resources orchestrated</div>
            </div>
          </div>
        </motion.div>
      </div>
    </div>
  );
}

