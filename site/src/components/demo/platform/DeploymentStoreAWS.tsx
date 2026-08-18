'use client';

import React from 'react';
import { motion } from 'framer-motion';
import { Network, Database, Server, Box, Shield, HardDrive } from 'lucide-react';

export default function DeploymentStoreAWS() {
  const awsComponents = [
    {
      name: 'AWS ALB',
      description: 'Application Load Balancer for high availability',
      icon: Network,
      color: 'from-white to-[#666]',
    },
    {
      name: 'AWS RDS',
      description: 'Managed relational database service',
      icon: Database,
      color: 'from-white to-[#666]',
    },
    {
      name: 'AWS EKS',
      description: 'Managed Kubernetes cluster on AWS',
      icon: Server,
      color: 'from-white to-[#666]',
    },
    {
      name: 'AWS ECS',
      description: 'Container orchestration with ECS',
      icon: Box,
      color: 'from-[#10b981] to-[#059669]',
    },
    {
      name: 'AWS VPC',
      description: 'Virtual Private Cloud networking',
      icon: Shield,
      color: 'from-[#999] to-[#666]',
    },
    {
      name: 'AWS S3',
      description: 'Scalable object storage service',
      icon: HardDrive,
      color: 'from-[#999] to-[#666]',
    },
  ];

  return (
    <div className="h-full overflow-y-auto bg-gradient-to-br from-[#151515] to-[#151515]">
      <div className="max-w-7xl mx-auto p-8 space-y-8">
        {/* Header */}
        <motion.div
          initial={{ opacity: 0, y: -20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.6 }}
          className="text-center"
        >
          <div className="inline-flex items-center justify-center w-20 h-20 bg-gradient-to-br from-white to-[#666] rounded-2xl mb-6 shadow-lg">
            <Server className="w-10 h-10 text-white" />
          </div>
          <h1 className="text-5xl font-black text-gray-900 mb-4">
            AWS Components
          </h1>
          <p className="text-xl text-gray-600 max-w-3xl mx-auto mb-6">
            100+ pre-built, production-ready AWS infrastructure components
          </p>
          <div className="inline-block bg-gradient-to-r from-[#999] to-red-600 text-white px-6 py-3 rounded-full">
            <p className="text-lg font-bold">
              100+ AWS components ready to deploy
            </p>
          </div>
        </motion.div>

        {/* Component Grid */}
        <div className="grid md:grid-cols-3 gap-6">
          {awsComponents.map((component, index) => {
            const Icon = component.icon;
            return (
              <motion.div
                key={component.name}
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: 0.1 * index + 0.3, duration: 0.5 }}
                className="bg-white rounded-2xl p-6 border-2 border-gray-200 demo-card-shadow hover:border-[#2a2a2a] hover:scale-105 transition-all duration-300"
              >
                <div className="flex flex-col items-center text-center gap-4">
                  <div className={`w-16 h-16 bg-gradient-to-br ${component.color} rounded-2xl flex items-center justify-center shadow-lg`}>
                    <Icon className="w-8 h-8 text-white" />
                  </div>
                  <div>
                    <h3 className="font-bold text-lg text-gray-900 mb-2">
                      {component.name}
                    </h3>
                    <p className="text-sm text-gray-600 leading-relaxed">
                      {component.description}
                    </p>
                  </div>
                </div>
              </motion.div>
            );
          })}
        </div>

        {/* Additional Categories */}
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.8, duration: 0.6 }}
          className="bg-white rounded-2xl p-8 border-2 border-[#2a2a2a]"
        >
          <h2 className="text-2xl font-bold text-gray-900 mb-6 text-center">
            And Many More AWS Resources
          </h2>
          <div className="grid md:grid-cols-4 gap-4 text-center">
            <div className="p-4">
              <h3 className="font-semibold text-gray-900 mb-2">Compute</h3>
              <p className="text-sm text-gray-600">EC2, Lambda, Fargate</p>
            </div>
            <div className="p-4">
              <h3 className="font-semibold text-gray-900 mb-2">Storage</h3>
              <p className="text-sm text-gray-600">EBS, EFS, ECR</p>
            </div>
            <div className="p-4">
              <h3 className="font-semibold text-gray-900 mb-2">Networking</h3>
              <p className="text-sm text-gray-600">Route53, Security Groups, NAT</p>
            </div>
            <div className="p-4">
              <h3 className="font-semibold text-gray-900 mb-2">Security</h3>
              <p className="text-sm text-gray-600">IAM, ACM, Secrets Manager</p>
            </div>
          </div>
        </motion.div>

        {/* Bottom Message */}
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 1.0, duration: 0.6 }}
          className="bg-gradient-to-r from-[#999] to-red-600 rounded-2xl p-6 text-white text-center"
        >
          <p className="text-lg font-medium">
            All components come with production-ready configurations and best practices built-in
          </p>
        </motion.div>
      </div>
    </div>
  );
}

