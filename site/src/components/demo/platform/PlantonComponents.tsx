'use client';

import React from 'react';
import { motion } from 'framer-motion';
import { Cloud, Server, Database, Network } from 'lucide-react';

export default function PlantonComponents() {
  const providerComponents = [
    {
      provider: 'AWS',
      icon: Cloud,
      color: 'from-white to-[#666]',
      components: ['ALB', 'NLB', 'RDS', 'DynamoDB', 'EKS', 'ECS', 'VPC', 'S3', 'ECR', 'Lambda', 'IAM', 'Security Groups'],
    },
    {
      provider: 'GCP',
      icon: Server,
      color: 'from-white to-[#666]',
      components: ['GKE', 'Cloud SQL', 'Load Balancing', 'VPC', 'Cloud Run', 'Cloud Storage', 'Cloud Functions', 'IAM'],
    },
    {
      provider: 'Azure',
      icon: Database,
      color: 'from-[#999] to-[#666]',
      components: ['AKS', 'Azure Database', 'Load Balancer', 'Virtual Network', 'Container Instances', 'Blob Storage', 'Functions'],
    },
    {
      provider: 'DigitalOcean',
      icon: Server,
      color: 'from-white to-[#666]',
      components: ['Droplets', 'Kubernetes', 'Managed Database', 'Spaces', 'Load Balancer', 'VPC'],
    },
    {
      provider: 'Auth0',
      icon: Database,
      color: 'from-[#10b981] to-[#059669]',
      components: ['Clients', 'Connections', 'Resource Servers', 'Event Streams'],
    },
    {
      provider: 'OpenFGA',
      icon: Network,
      color: 'from-white to-[#666]',
      components: ['Stores', 'Authorization Models', 'Relationship Tuples'],
    },
    {
      provider: 'And More',
      icon: Cloud,
      color: 'from-white to-[#666]',
      components: ['Cloudflare', 'DigitalOcean', 'Kubernetes'],
    },
  ];

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
          <h1 className="text-5xl font-black text-gray-900 mb-4">
            500+ Infrastructure Components
          </h1>
          <p className="text-xl text-gray-600 max-w-3xl mx-auto mb-6">
            Pre-defined component specifications covering all major cloud providers and platforms
          </p>
          <div className="inline-block bg-gradient-to-r from-white to-[#666] text-white px-6 py-3 rounded-full">
            <p className="text-lg font-bold">
              Consistent API across all providers
            </p>
          </div>
        </motion.div>

        {/* Provider Component Grid */}
        <div className="grid md:grid-cols-3 gap-6">
          {providerComponents.map((provider, index) => {
            const Icon = provider.icon;
            return (
              <motion.div
                key={provider.provider}
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: 0.1 * index + 0.3, duration: 0.5 }}
                className="bg-white rounded-2xl border-2 border-gray-200 overflow-hidden demo-card-shadow"
              >
                {/* Provider Header */}
                <div className={`bg-gradient-to-r ${provider.color} px-6 py-4`}>
                  <div className="flex items-center gap-3">
                    <div className="w-10 h-10 bg-white/20 rounded-lg flex items-center justify-center">
                      <Icon className="w-5 h-5 text-white" />
                    </div>
                    <h3 className="font-bold text-lg text-white">{provider.provider}</h3>
                  </div>
                </div>

                {/* Components List */}
                <div className="p-6">
                  <div className="flex flex-wrap gap-2">
                    {provider.components.map((component, i) => (
                      <span
                        key={i}
                        className="px-3 py-1 bg-gray-100 text-gray-700 text-xs font-medium rounded-full border border-gray-200"
                      >
                        {component}
                      </span>
                    ))}
                  </div>
                </div>
              </motion.div>
            );
          })}
        </div>

        {/* Features */}
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.8, duration: 0.6 }}
          className="bg-gradient-to-r from-[#151515] to-[#151515] rounded-2xl p-8 border-2 border-[#2a2a2a]"
        >
          <h2 className="text-2xl font-bold text-gray-900 mb-6 text-center">
            What Makes These Components Special
          </h2>
          <div className="grid md:grid-cols-3 gap-6">
            <div className="text-center">
              <div className="w-12 h-12 bg-[#111] rounded-xl flex items-center justify-center mx-auto mb-3">
                <Server className="w-6 h-6 text-white" />
              </div>
              <h3 className="font-bold text-gray-900 mb-2">Type-Safe Specs</h3>
              <p className="text-sm text-gray-700">
                YAML validated against schemas for correctness
              </p>
            </div>
            <div className="text-center">
              <div className="w-12 h-12 bg-[#111] rounded-xl flex items-center justify-center mx-auto mb-3">
                <Network className="w-6 h-6 text-white" />
              </div>
              <h3 className="font-bold text-gray-900 mb-2">Cross-References</h3>
              <p className="text-sm text-gray-700">
                Link resources together with valueFrom
              </p>
            </div>
            <div className="text-center">
              <div className="w-12 h-12 bg-[#333] rounded-xl flex items-center justify-center mx-auto mb-3">
                <Database className="w-6 h-6 text-white" />
              </div>
              <h3 className="font-bold text-gray-900 mb-2">Production Ready</h3>
              <p className="text-sm text-gray-700">
                Best practices and sensible defaults
              </p>
            </div>
          </div>
        </motion.div>

        {/* Bottom CTA */}
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 1.0, duration: 0.6 }}
          className="bg-gradient-to-r from-white to-[#666] rounded-2xl p-6 text-white text-center"
        >
          <p className="text-lg font-medium">
            All component definitions are open source and available on GitHub
          </p>
        </motion.div>
      </div>
    </div>
  );
}
