'use client';

import React from 'react';
import { motion } from 'framer-motion';
import { DollarSign, Code, Users } from 'lucide-react';

export default function IntroProblemSmallProduct() {
  return (
    <div className="h-full flex flex-col justify-center">
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.6 }}
        className="text-center max-w-4xl mx-auto px-8"
      >
        <h1 className="text-5xl md:text-6xl lg:text-7xl font-black text-gray-900 mb-6 tracking-tight">
          Small teams can&apos;t afford <br />
          <span className="text-red-600">dedicated DevOps.</span>
        </h1>
        
        <p className="text-xl md:text-2xl text-gray-600 mb-12 max-w-3xl mx-auto leading-relaxed">
          Developers are stuck managing infrastructure instead of building features, and hiring a DevOps team is out of reach.
        </p>

        <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-12">
          <motion.div
            initial={{ opacity: 0, scale: 0.9 }}
            animate={{ opacity: 1, scale: 1 }}
            transition={{ delay: 0.2 }}
            className="bg-red-50 border border-red-200 rounded-xl p-6"
          >
            <DollarSign className="w-10 h-10 text-red-600 mb-3 mx-auto" />
            <h3 className="font-bold text-gray-900 mb-2">Scarce</h3>
            <p className="text-gray-600 text-sm">Senior ops expertise is hard to hire</p>
          </motion.div>

          <motion.div
            initial={{ opacity: 0, scale: 0.9 }}
            animate={{ opacity: 1, scale: 1 }}
            transition={{ delay: 0.3 }}
            className="bg-red-50 border border-red-200 rounded-xl p-6"
          >
            <Code className="w-10 h-10 text-red-600 mb-3 mx-auto" />
            <h3 className="font-bold text-gray-900 mb-2">40% Time Waste</h3>
            <p className="text-gray-600 text-sm">Developers spending time on ops instead of features</p>
          </motion.div>

          <motion.div
            initial={{ opacity: 0, scale: 0.9 }}
            animate={{ opacity: 1, scale: 1 }}
            transition={{ delay: 0.4 }}
            className="bg-red-50 border border-red-200 rounded-xl p-6"
          >
            <Users className="w-10 h-10 text-red-600 mb-3 mx-auto" />
            <h3 className="font-bold text-gray-900 mb-2">Scaling Blocked</h3>
            <p className="text-gray-600 text-sm">Infrastructure becomes bottleneck as team grows</p>
          </motion.div>
        </div>

        <motion.div
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ delay: 0.6 }}
          className="text-gray-500 italic"
        >
          <p className="text-lg">
            &ldquo;Our developers should ship features, not debug Terraform.&rdquo;
          </p>
          <p className="text-sm mt-2">— Every Product Startup CTO</p>
        </motion.div>
      </motion.div>
    </div>
  );
}

