'use client';

import React from 'react';
import { motion } from 'framer-motion';
import { Cloud, AlertCircle, DollarSign, Clock } from 'lucide-react';

export default function IntroProblem() {
  return (
    <div className="h-full flex flex-col justify-center">
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.6 }}
        className="text-center max-w-4xl mx-auto px-8"
      >
        <h1 className="text-5xl md:text-6xl lg:text-7xl font-black text-gray-900 mb-6 tracking-tight">
          Multi-cloud deployments <br />
          <span className="text-red-600">are broken.</span>
        </h1>
        
        <p className="text-xl md:text-2xl text-gray-600 mb-12 max-w-3xl mx-auto leading-relaxed">
          Teams waste months juggling tools, fixing errors, and overspending on misconfigurations.
        </p>

        <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-12">
          <motion.div
            initial={{ opacity: 0, scale: 0.9 }}
            animate={{ opacity: 1, scale: 1 }}
            transition={{ delay: 0.2 }}
            className="bg-red-50 border border-red-200 rounded-xl p-6"
          >
            <Clock className="w-10 h-10 text-red-600 mb-3 mx-auto" />
            <h3 className="font-bold text-gray-900 mb-2">6+ Months</h3>
            <p className="text-gray-600 text-sm">Average setup time for multi-cloud infrastructure</p>
          </motion.div>

          <motion.div
            initial={{ opacity: 0, scale: 0.9 }}
            animate={{ opacity: 1, scale: 1 }}
            transition={{ delay: 0.3 }}
            className="bg-red-50 border border-red-200 rounded-xl p-6"
          >
            <AlertCircle className="w-10 h-10 text-red-600 mb-3 mx-auto" />
            <h3 className="font-bold text-gray-900 mb-2">70% Errors</h3>
            <p className="text-gray-600 text-sm">Due to manual configuration and tool fragmentation</p>
          </motion.div>

          <motion.div
            initial={{ opacity: 0, scale: 0.9 }}
            animate={{ opacity: 1, scale: 1 }}
            transition={{ delay: 0.4 }}
            className="bg-red-50 border border-red-200 rounded-xl p-6"
          >
            <DollarSign className="w-10 h-10 text-red-600 mb-3 mx-auto" />
            <h3 className="font-bold text-gray-900 mb-2">40% Overspend</h3>
            <p className="text-gray-600 text-sm">From inefficient resource allocation and sprawl</p>
          </motion.div>
        </div>

        <motion.div
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ delay: 0.6 }}
          className="flex justify-center gap-8"
        >
          <div className="opacity-40 grayscale">
            <Cloud className="h-12 w-12 text-[#a0a0a0]" />
          </div>
          <div className="opacity-40 grayscale">
            <Cloud className="h-12 w-12 text-white" />
          </div>
          <div className="opacity-40 grayscale">
            <Cloud className="h-12 w-12 text-white" />
          </div>
        </motion.div>
      </motion.div>
    </div>
  );
}
