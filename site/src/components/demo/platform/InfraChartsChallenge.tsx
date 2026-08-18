'use client';

import React from 'react';
import { motion } from 'framer-motion';
import { Lightbulb } from 'lucide-react';

export default function InfraChartsChallenge() {
  return (
    <div className="h-full overflow-y-auto flex items-center justify-center bg-gradient-to-br from-red-50 via-[#151515] to-[#151515]">
      <div className="max-w-5xl mx-auto p-8">
        <motion.div
          initial={{ opacity: 0, scale: 0.95 }}
          animate={{ opacity: 1, scale: 1 }}
          transition={{ duration: 0.6 }}
          className="bg-white rounded-3xl border-2 border-red-200 p-12 demo-card-shadow"
        >
          {/* Header */}
          <div className="text-center mb-8">
            <div className="inline-flex items-center justify-center w-20 h-20 bg-red-600 rounded-2xl mb-6">
              <Lightbulb className="w-10 h-10 text-white" />
            </div>
            <h1 className="text-5xl font-black text-gray-900 mb-4">
              The Challenge
            </h1>
            <p className="text-xl text-gray-600 max-w-2xl mx-auto">
              Production environments aren&apos;t just one resource
            </p>
          </div>

          {/* Main Content */}
          <div className="space-y-6">
            <p className="text-2xl text-gray-700 text-center font-medium">
              An ECS environment needs:
            </p>
            
            <div className="grid grid-cols-2 md:grid-cols-3 gap-4 mt-8">
              {[
                'VPC with subnets',
                'Security groups',
                'Application Load Balancer',
                'SSL certificates',
                'Route53 DNS',
                'ECS cluster',
                'IAM roles',
                'ECR repositories',
                'ECS services'
              ].map((item, i) => (
                <motion.div
                  key={i}
                  initial={{ opacity: 0, y: 20 }}
                  animate={{ opacity: 1, y: 0 }}
                  transition={{ delay: 0.1 * i, duration: 0.4 }}
                  className="bg-gradient-to-br from-red-50 to-[#151515] rounded-xl px-6 py-4 text-center font-semibold text-gray-800 border-2 border-red-200"
                >
                  {item}
                </motion.div>
              ))}
            </div>

            <motion.div
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: 1.0, duration: 0.6 }}
              className="mt-10 pt-8 border-t-2 border-red-200"
            >
              <p className="text-2xl font-bold text-red-900 text-center leading-relaxed">
                Manually deploying and connecting these is complex, error-prone, and time-consuming.
              </p>
            </motion.div>
          </div>
        </motion.div>
      </div>
    </div>
  );
}

