'use client';

import React from 'react';
import { motion } from 'framer-motion';
import { Award, Clock, Zap, CheckCircle, Sparkles } from 'lucide-react';

export default function ServiceSuccessStory() {
  return (
    <div className="h-full flex flex-col">
      {/* Header */}
      <div className="px-8 py-6 border-b border-gray-200">
        <h2 className="text-2xl font-bold text-gray-900">Customer Success Story</h2>
        <p className="text-gray-600 mt-1">
          Real results from IT consulting companies using Planton
        </p>
      </div>

      <div className="flex-1 overflow-auto bg-gradient-to-br from-[#151515] via-[#151515] to-[#151515]">
        <div className="max-w-5xl mx-auto p-8">
          {/* Hero section */}
          <motion.div
            initial={{ opacity: 0, y: 30 }}
            animate={{ opacity: 1, y: 0 }}
            className="text-center mb-12"
          >
            <div className="inline-flex items-center justify-center w-20 h-20 bg-gradient-to-br from-white to-[#666] rounded-2xl mb-6 shadow-xl">
              <Sparkles className="w-10 h-10 text-white" />
            </div>
            <h1 className="text-5xl font-extrabold text-gray-900 mb-4">
              From Code to Production
            </h1>
            <p className="text-3xl text-white font-bold mb-6">
              In Under 1 Hour
            </p>
          </motion.div>

          {/* Main success card */}
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.2 }}
            className="bg-gradient-to-r from-white to-[#666] rounded-2xl p-8 text-white shadow-2xl mb-8"
          >
            <div className="flex items-start gap-4 mb-8">
              <div className="w-12 h-12 bg-white/20 backdrop-blur-sm rounded-xl flex items-center justify-center flex-shrink-0">
                <Award className="w-6 h-6 text-white" />
              </div>
              <div>
                <h3 className="text-2xl font-bold mb-2">IT Consulting Company Success</h3>
                <p className="text-gray-300">
                  Deploying infrastructure for 3 different clients with zero duplication
                </p>
              </div>
            </div>

            <div className="grid md:grid-cols-2 gap-6 mb-8">
              <motion.div
                initial={{ opacity: 0, x: -20 }}
                animate={{ opacity: 1, x: 0 }}
                transition={{ delay: 0.3 }}
                className="bg-white/10 backdrop-blur-sm rounded-xl p-6 border border-white/20"
              >
                <div className="flex items-center gap-3 mb-3">
                  <Clock className="w-8 h-8 text-white" />
                  <div className="text-4xl font-bold">20 min</div>
                </div>
                <div className="text-gray-300">
                  Full infrastructure deployment with EKS cluster, databases, and networking
                </div>
              </motion.div>

              <motion.div
                initial={{ opacity: 0, x: 20 }}
                animate={{ opacity: 1, x: 0 }}
                transition={{ delay: 0.4 }}
                className="bg-white/10 backdrop-blur-sm rounded-xl p-6 border border-white/20"
              >
                <div className="flex items-center gap-3 mb-3">
                  <Zap className="w-8 h-8 text-white" />
                  <div className="text-4xl font-bold">&lt; 1 hour</div>
                </div>
                <div className="text-gray-300">
                  From code commit to live URL in client&apos;s domain with SSL certificate
                </div>
              </motion.div>
            </div>

            <div className="pt-6 border-t border-white/20">
              <p className="text-lg text-gray-300 italic">
                &ldquo;Every single time they need to duplicate terraform code for a new client, Planton saves hours of work. 
                The same standard approach across all projects means faster onboarding and fewer errors.&rdquo;
              </p>
            </div>
          </motion.div>

          {/* Key benefits grid */}
          <div className="grid md:grid-cols-3 gap-6 mb-8">
            <motion.div
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: 0.5 }}
              className="bg-white rounded-2xl p-6 shadow-lg border border-gray-200"
            >
              <div className="flex items-center gap-3 mb-4">
                <div className="w-12 h-12 bg-gradient-to-br from-[#10b981] to-[#059669] rounded-xl flex items-center justify-center">
                  <CheckCircle className="w-6 h-6 text-white" />
                </div>
                <h3 className="font-bold text-gray-900">No Code Duplication</h3>
              </div>
              <p className="text-gray-600 text-sm">
                Deploy the same infrastructure stack for multiple clients without copying terraform files
              </p>
            </motion.div>

            <motion.div
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: 0.6 }}
              className="bg-white rounded-2xl p-6 shadow-lg border border-gray-200"
            >
              <div className="flex items-center gap-3 mb-4">
                <div className="w-12 h-12 bg-gradient-to-br from-white to-[#666] rounded-xl flex items-center justify-center">
                  <Zap className="w-6 h-6 text-white" />
                </div>
                <h3 className="font-bold text-gray-900">Rapid Deployment</h3>
              </div>
              <p className="text-gray-600 text-sm">
                New client environments ready in minutes, not days. Complete infrastructure in 20 minutes
              </p>
            </motion.div>

            <motion.div
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: 0.7 }}
              className="bg-white rounded-2xl p-6 shadow-lg border border-gray-200"
            >
              <div className="flex items-center gap-3 mb-4">
                <div className="w-12 h-12 bg-gradient-to-br from-white to-[#666] rounded-xl flex items-center justify-center">
                  <Award className="w-6 h-6 text-white" />
                </div>
                <h3 className="font-bold text-gray-900">Standard Approach</h3>
              </div>
              <p className="text-gray-600 text-sm">
                Same proven patterns across all projects means predictable results and easier team onboarding
              </p>
            </motion.div>
          </div>

          {/* Final CTA */}
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.8 }}
            className="bg-white rounded-2xl p-8 shadow-xl border border-gray-200 text-center"
          >
            <h3 className="text-2xl font-bold text-gray-900 mb-4">
              Ready to Transform Your DevOps?
            </h3>
            <p className="text-gray-600 mb-6 max-w-2xl mx-auto">
              Join IT consulting companies and product teams who have eliminated DevOps overhead and accelerated their delivery with Planton.
            </p>
            <div className="flex flex-wrap items-center justify-center gap-4">
              <div className="flex items-center gap-2 text-sm text-gray-600">
                <CheckCircle className="w-4 h-4 text-[#10b981]" />
                <span>Infrastructure in 20 minutes</span>
              </div>
              <div className="flex items-center gap-2 text-sm text-gray-600">
                <CheckCircle className="w-4 h-4 text-[#10b981]" />
                <span>Services deployed in &lt; 1 hour</span>
              </div>
              <div className="flex items-center gap-2 text-sm text-gray-600">
                <CheckCircle className="w-4 h-4 text-[#10b981]" />
                <span>Zero code duplication</span>
              </div>
            </div>
          </motion.div>
        </div>
      </div>
    </div>
  );
}

