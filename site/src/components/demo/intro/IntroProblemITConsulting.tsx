'use client';

import React from 'react';
import { motion } from 'framer-motion';
import { Clock, TrendingDown, AlertTriangle } from 'lucide-react';

export default function IntroProblemITConsulting() {
  return (
    <div className="h-full flex flex-col justify-center">
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.6 }}
        className="text-center max-w-4xl mx-auto px-8"
      >
        <h1 className="text-5xl md:text-6xl lg:text-7xl font-black text-gray-900 mb-6 tracking-tight">
          Client projects delayed by <br />
          <span className="text-red-600">infrastructure setup.</span>
        </h1>
        
        <p className="text-xl md:text-2xl text-gray-600 mb-12 max-w-3xl mx-auto leading-relaxed">
          Every new client engagement means weeks of repetitive infrastructure work, delaying time-to-value and losing competitive deals.
        </p>

        <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-12">
          <motion.div
            initial={{ opacity: 0, scale: 0.9 }}
            animate={{ opacity: 1, scale: 1 }}
            transition={{ delay: 0.2 }}
            className="bg-red-50 border border-red-200 rounded-xl p-6"
          >
            <Clock className="w-10 h-10 text-red-600 mb-3 mx-auto" />
            <h3 className="font-bold text-gray-900 mb-2">3+ Weeks</h3>
            <p className="text-gray-600 text-sm">Per client project setup time for infrastructure</p>
          </motion.div>

          <motion.div
            initial={{ opacity: 0, scale: 0.9 }}
            animate={{ opacity: 1, scale: 1 }}
            transition={{ delay: 0.3 }}
            className="bg-red-50 border border-red-200 rounded-xl p-6"
          >
            <TrendingDown className="w-10 h-10 text-red-600 mb-3 mx-auto" />
            <h3 className="font-bold text-gray-900 mb-2">Lost Deals</h3>
            <p className="text-gray-600 text-sm">Clients choosing competitors with faster delivery</p>
          </motion.div>

          <motion.div
            initial={{ opacity: 0, scale: 0.9 }}
            animate={{ opacity: 1, scale: 1 }}
            transition={{ delay: 0.4 }}
            className="bg-red-50 border border-red-200 rounded-xl p-6"
          >
            <AlertTriangle className="w-10 h-10 text-red-600 mb-3 mx-auto" />
            <h3 className="font-bold text-gray-900 mb-2">Onboarding Delays</h3>
            <p className="text-gray-600 text-sm">Repetitive DevOps work blocking billable hours</p>
          </motion.div>
        </div>

        <motion.div
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ delay: 0.6 }}
          className="text-gray-500 italic"
        >
          <p className="text-lg">
            &ldquo;We need infrastructure ready on day one, not week three.&rdquo;
          </p>
          <p className="text-sm mt-2">— Every IT Consulting Client</p>
        </motion.div>
      </motion.div>
    </div>
  );
}

