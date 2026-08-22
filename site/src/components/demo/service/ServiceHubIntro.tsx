'use client';

import React from 'react';
import { motion } from 'framer-motion';
import { Zap, GitBranch, Rocket, CheckCircle } from 'lucide-react';

export default function ServiceHubIntro() {
  return (
    <div className="h-full flex flex-col">
      {/* Header */}
      <div className="px-8 py-6 border-b border-gray-200 bg-gradient-to-r from-[#151515] to-[#151515]">
        <motion.div
          initial={{ opacity: 0, y: -20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.6 }}
        >
          <h2 className="text-3xl font-bold text-gray-900 mb-2">Service Hub</h2>
          <p className="text-xl text-gray-600 mb-2">
            Vercel for Backend, In Your Own Cloud
          </p>
          <div className="inline-block bg-gradient-to-r from-[#10b981] to-[#059669] text-white px-5 py-2 rounded-full text-sm font-bold">
            ✨ No Dockerfile Required
          </div>
        </motion.div>
      </div>

      <div className="flex-1 overflow-auto bg-gradient-to-br from-[#151515] via-[#111] to-[#151515]">
        <div className="max-w-5xl mx-auto p-12">
          {/* Hero section */}
          <motion.div
            initial={{ opacity: 0, y: 30 }}
            animate={{ opacity: 1, y: 0 }}
            className="text-center mb-12"
          >
            <div className="inline-flex items-center justify-center w-20 h-20 bg-gradient-to-br from-white to-[#666] rounded-2xl mb-6 shadow-xl">
              <Rocket className="w-10 h-10 text-white" />
            </div>
            <h1 className="text-5xl font-extrabold text-gray-900 mb-4">
              Service Hub
            </h1>
            <p className="text-2xl text-gray-600 font-medium mb-4">
              Vercel for Backend, In Your Own Cloud
            </p>
            <div className="inline-block bg-gradient-to-r from-[#10b981] to-[#059669] text-white px-8 py-4 rounded-2xl mb-4 shadow-xl">
              <p className="text-2xl font-black">
                No Dockerfile Required
              </p>
            </div>
            <p className="text-lg text-gray-500 max-w-2xl mx-auto">
              Drop-in replacement for GitHub Actions and GitLab pipelines. Deploy backend services with zero DevOps overhead.
            </p>
          </motion.div>

          {/* Key benefits grid */}
          <div className="grid md:grid-cols-2 gap-6 mb-12">
            <motion.div
              initial={{ opacity: 0, x: -20 }}
              animate={{ opacity: 1, x: 0 }}
              transition={{ delay: 0.1 }}
              className="bg-white rounded-2xl p-8 shadow-lg border border-gray-200"
            >
              <div className="flex items-center gap-3 mb-4">
                <div className="w-12 h-12 bg-gradient-to-br from-[#10b981] to-[#059669] rounded-xl flex items-center justify-center">
                  <CheckCircle className="w-6 h-6 text-white" />
                </div>
                <h3 className="font-bold text-gray-900 text-lg">No Jenkins Required</h3>
              </div>
              <p className="text-gray-600">
                Built-in CI/CD pipeline handles everything from code commit to production deployment. No complex pipeline configuration needed.
              </p>
            </motion.div>

            <motion.div
              initial={{ opacity: 0, x: 20 }}
              animate={{ opacity: 1, x: 0 }}
              transition={{ delay: 0.2 }}
              className="bg-white rounded-2xl p-8 shadow-lg border border-gray-200"
            >
              <div className="flex items-center gap-3 mb-4">
                <div className="w-12 h-12 bg-gradient-to-br from-white to-[#666] rounded-xl flex items-center justify-center">
                  <GitBranch className="w-6 h-6 text-white" />
                </div>
                <h3 className="font-bold text-gray-900 text-lg">No GitHub Actions</h3>
              </div>
              <p className="text-gray-600">
                Platform provides fully managed build and deployment infrastructure. Connect your repo and you&apos;re done.
              </p>
            </motion.div>

            <motion.div
              initial={{ opacity: 0, x: -20 }}
              animate={{ opacity: 1, x: 0 }}
              transition={{ delay: 0.3 }}
              className="bg-white rounded-2xl p-8 shadow-lg border-2 border-[#10b981]/30"
            >
              <div className="flex items-center gap-3 mb-4">
                <div className="w-12 h-12 bg-gradient-to-br from-[#10b981] to-[#059669] rounded-xl flex items-center justify-center">
                  <Zap className="w-6 h-6 text-white" />
                </div>
                <h3 className="font-bold text-gray-900 text-lg">No Dockerfile Required</h3>
              </div>
              <p className="text-gray-600">
                Cloud-native BuildPacks automatically detect your tech stack and create optimized container images. Zero Dockerfile knowledge needed.
              </p>
            </motion.div>

            <motion.div
              initial={{ opacity: 0, x: 20 }}
              animate={{ opacity: 1, x: 0 }}
              transition={{ delay: 0.4 }}
              className="bg-white rounded-2xl p-8 shadow-lg border border-gray-200"
            >
              <div className="flex items-center gap-3 mb-4">
                <div className="w-12 h-12 bg-gradient-to-br from-white to-[#666] rounded-xl flex items-center justify-center">
                  <Zap className="w-6 h-6 text-white" />
                </div>
                <h3 className="font-bold text-gray-900 text-lg">Deploy in Minutes</h3>
              </div>
              <p className="text-gray-600">
                From code commit to live URL in under an hour. 20 minutes for infrastructure, then instant deployments for every commit.
              </p>
            </motion.div>
          </div>

          {/* Comparison section */}
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.5 }}
            className="bg-gradient-to-r from-white to-[#666] rounded-2xl p-8 shadow-2xl text-white"
          >
            <h3 className="text-2xl font-bold mb-6 text-center">
              Two Ways to Deploy
            </h3>
            <div className="grid md:grid-cols-2 gap-8">
              <div className="bg-white/10 backdrop-blur-sm rounded-xl p-6 border border-white/20">
                <h4 className="font-bold text-lg mb-3">Bring Your Dockerfile</h4>
                <ul className="space-y-2 text-sm text-white/90">
                  <li>✓ Use your existing Dockerfile</li>
                  <li>✓ Full control over build process</li>
                  <li>✓ Custom base images</li>
                  <li>✓ Multi-stage builds supported</li>
                </ul>
              </div>
              <div className="bg-white/10 backdrop-blur-sm rounded-xl p-6 border-2 border-[#10b981]">
                <h4 className="font-bold text-lg mb-1">No Dockerfile Required</h4>
                <p className="text-xs text-[#10b981] mb-3">Recommended for most teams</p>
                <ul className="space-y-2 text-sm text-white/90">
                  <li>✓ Zero Dockerfile knowledge needed</li>
                  <li>✓ Auto-detect tech stack</li>
                  <li>✓ Support for Java, Python, Node.js, Go...</li>
                  <li>✓ Best practices built-in</li>
                  <li>✓ Automatic security updates</li>
                </ul>
              </div>
            </div>
          </motion.div>

          {/* CTA */}
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.6 }}
            className="mt-12 text-center"
          >
            <p className="text-gray-600 mb-6">
              Let&apos;s see how to deploy a backend service in action
            </p>
          </motion.div>
        </div>
      </div>
    </div>
  );
}

