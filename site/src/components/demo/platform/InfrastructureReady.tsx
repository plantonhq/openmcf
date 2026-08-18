'use client';

import React from 'react';
import { motion } from 'framer-motion';
import { CheckCircle2, Rocket, ArrowRight, Server, Network, Shield } from 'lucide-react';

export default function InfrastructureReady() {
  return (
    <div className="h-full flex flex-col justify-center">
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.6 }}
        className="text-center max-w-5xl mx-auto px-8"
      >
        {/* Success Badge */}
        <motion.div
          initial={{ scale: 0 }}
          animate={{ scale: 1 }}
          transition={{ delay: 0.2, type: "spring", stiffness: 200 }}
          className="inline-flex items-center justify-center w-24 h-24 bg-gradient-to-r from-[#10b981] to-[#666] rounded-full mb-8 shadow-xl"
        >
          <CheckCircle2 className="w-12 h-12 text-white" />
        </motion.div>

        {/* Main Message */}
        <motion.h1
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.4 }}
          className="text-5xl md:text-6xl lg:text-7xl font-black text-gray-900 mb-6 tracking-tight"
        >
          Infrastructure is <br />
          <span className="text-transparent bg-clip-text bg-gradient-to-r from-[#10b981] to-[#059669]">
            Production Ready
          </span>
        </motion.h1>

        <motion.p
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.6 }}
          className="text-xl md:text-2xl text-gray-600 mb-12 max-w-3xl mx-auto leading-relaxed"
        >
          Your AWS environment is fully configured and ready. Developers can now start deploying services.
        </motion.p>

        {/* What's Ready */}
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.8 }}
          className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-12 max-w-4xl mx-auto"
        >
          <div className="bg-white border-2 border-[#2a2a2a] rounded-2xl p-6 demo-card-shadow">
            <div className="w-12 h-12 bg-gradient-to-br from-[#10b981] to-[#666] rounded-xl flex items-center justify-center mx-auto mb-4">
              <Network className="w-6 h-6 text-white" />
            </div>
            <h3 className="font-bold text-lg text-gray-900 mb-2">Networking</h3>
            <p className="text-sm text-gray-600">
              VPC, subnets, security groups, and load balancer configured
            </p>
          </div>

          <div className="bg-white border-2 border-[#2a2a2a] rounded-2xl p-6 demo-card-shadow">
            <div className="w-12 h-12 bg-gradient-to-br from-[#10b981] to-[#666] rounded-xl flex items-center justify-center mx-auto mb-4">
              <Server className="w-6 h-6 text-white" />
            </div>
            <h3 className="font-bold text-lg text-gray-900 mb-2">Compute</h3>
            <p className="text-sm text-gray-600">
              ECS cluster ready to run containerized workloads
            </p>
          </div>

          <div className="bg-white border-2 border-[#2a2a2a] rounded-2xl p-6 demo-card-shadow">
            <div className="w-12 h-12 bg-gradient-to-br from-[#10b981] to-[#666] rounded-xl flex items-center justify-center mx-auto mb-4">
              <Shield className="w-6 h-6 text-white" />
            </div>
            <h3 className="font-bold text-lg text-gray-900 mb-2">Security</h3>
            <p className="text-sm text-gray-600">
              SSL certificates, DNS, and IAM roles configured
            </p>
          </div>
        </motion.div>

        {/* Stats */}
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 1.0 }}
          className="bg-gradient-to-r from-[#151515] to-[#151515] border-2 border-[#2a2a2a] rounded-2xl p-8 mb-12 max-w-4xl mx-auto"
        >
          <div className="grid grid-cols-3 gap-8 mb-6">
            <div>
              <div className="text-4xl font-black text-[#10b981] mb-2">9</div>
              <div className="text-sm text-gray-600">Resources Deployed</div>
            </div>
            <div>
              <div className="text-4xl font-black text-[#10b981] mb-2">~20</div>
              <div className="text-sm text-gray-600">Minutes Total Time</div>
            </div>
            <div>
              <div className="text-4xl font-black text-[#10b981] mb-2">✓</div>
              <div className="text-sm text-gray-600">Production Ready</div>
            </div>
          </div>
          <p className="text-gray-700 font-medium">
            All underlying infrastructure is fully configured. Your environment is ready for service deployments.
          </p>
        </motion.div>

        {/* Next Step - Segue to Service Hub */}
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 1.2 }}
          className="bg-gradient-to-r from-white to-[#666] rounded-2xl p-8 text-white max-w-4xl mx-auto"
        >
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-4">
              <div className="w-14 h-14 bg-white/20 rounded-xl flex items-center justify-center backdrop-blur-sm">
                <Rocket className="w-7 h-7 text-white" />
              </div>
              <div className="text-left">
                <h3 className="text-2xl font-black mb-1">Next: Deploy Services</h3>
                <p className="text-white/90 text-sm">
                  Infrastructure is ready. Let&apos;s see how to deploy your applications to this environment.
                </p>
              </div>
            </div>
            <ArrowRight className="w-8 h-8 flex-shrink-0" />
          </div>
        </motion.div>

        {/* Additional Context */}
        <motion.div
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ delay: 1.4 }}
          className="mt-8"
        >
          <p className="text-gray-500 text-sm italic">
            This is where Infrastructure Hub hands off to Service Hub—the next part of the platform.
          </p>
        </motion.div>
      </motion.div>
    </div>
  );
}

