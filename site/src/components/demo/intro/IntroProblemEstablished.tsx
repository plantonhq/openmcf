'use client';

import React from 'react';
import { motion } from 'framer-motion';
import { ShieldAlert, Eye, AlertCircle } from 'lucide-react';

export default function IntroProblemEstablished() {
  return (
    <div className="h-full flex flex-col justify-center">
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.6 }}
        className="text-center max-w-4xl mx-auto px-8"
      >
        <h1 className="text-5xl md:text-6xl lg:text-7xl font-black text-gray-900 mb-6 tracking-tight">
          Infrastructure sprawl <br />
          <span className="text-red-600">creates blind spots.</span>
        </h1>
        
        <p className="text-xl md:text-2xl text-gray-600 mb-12 max-w-3xl mx-auto leading-relaxed">
          Multiple teams deploying infrastructure independently leads to compliance gaps, security risks, and zero visibility.
        </p>

        <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-12">
          <motion.div
            initial={{ opacity: 0, scale: 0.9 }}
            animate={{ opacity: 1, scale: 1 }}
            transition={{ delay: 0.2 }}
            className="bg-red-50 border border-red-200 rounded-xl p-6"
          >
            <ShieldAlert className="w-10 h-10 text-red-600 mb-3 mx-auto" />
            <h3 className="font-bold text-gray-900 mb-2">73% Have Gaps</h3>
            <p className="text-gray-600 text-sm">Organizations with compliance violations</p>
          </motion.div>

          <motion.div
            initial={{ opacity: 0, scale: 0.9 }}
            animate={{ opacity: 1, scale: 1 }}
            transition={{ delay: 0.3 }}
            className="bg-red-50 border border-red-200 rounded-xl p-6"
          >
            <AlertCircle className="w-10 h-10 text-red-600 mb-3 mx-auto" />
            <h3 className="font-bold text-gray-900 mb-2">28 Incidents/Year</h3>
            <p className="text-gray-600 text-sm">Average security misconfigurations per company</p>
          </motion.div>

          <motion.div
            initial={{ opacity: 0, scale: 0.9 }}
            animate={{ opacity: 1, scale: 1 }}
            transition={{ delay: 0.4 }}
            className="bg-red-50 border border-red-200 rounded-xl p-6"
          >
            <Eye className="w-10 h-10 text-red-600 mb-3 mx-auto" />
            <h3 className="font-bold text-gray-900 mb-2">Zero Visibility</h3>
            <p className="text-gray-600 text-sm">No centralized view of infrastructure state</p>
          </motion.div>
        </div>

        <motion.div
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ delay: 0.6 }}
          className="text-gray-500 italic"
        >
          <p className="text-lg">
            &ldquo;We need golden paths and guardrails, not chaos.&rdquo;
          </p>
          <p className="text-sm mt-2">— Every Platform Engineering Lead</p>
        </motion.div>
      </motion.div>
    </div>
  );
}

