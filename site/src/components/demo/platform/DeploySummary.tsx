'use client';

import React, { useState } from 'react';
import { motion } from 'framer-motion';
import {
  CheckCircle,
  Clock,
  DollarSign,
  Zap,
  AlertTriangle,
  Server,
} from 'lucide-react';
import AwsVpcIcon from '../../../assets/icons/aws-vpc.svg';
import AwsEksClusterIcon from '../../../assets/icons/aws-eks-cluster.svg';
import PostgresKubernetesIcon from '../../../assets/icons/postgres-kubernetes.svg';
import RedisKubernetesIcon from '../../../assets/icons/redis-kubernetes.svg';
import CloudflareR2BucketIcon from '../../../assets/icons/cloudflare-r2-bucket.svg';
import AwsCloudfrontIcon from '../../../assets/icons/aws-cloudfront.svg';
import DeployResourceItem from './DeployResourceItem';

interface DeploySummaryProps {
  onDeploy: () => void;
}

export default function DeploySummary({ onDeploy }: DeploySummaryProps) {
  const [isReviewing, setIsReviewing] = useState(false);

  const deploymentItems = [
    {
      id: 'vpc-1',
      name: 'VPC',
      icon: AwsVpcIcon,
      color: 'from-white to-[#666]',
      config: { region: 'us-east-1', cidr: '10.0.0.0/16' },
      cost: '$0/month',
    },
    {
      id: 'eks-1',
      name: 'EKS Cluster',
      icon: AwsEksClusterIcon,
      color: 'from-white to-[#666]',
      config: { region: 'us-east-1', size: 'medium', nodes: 3 },
      cost: '$216/month',
    },
    {
      id: 'postgres-1',
      name: 'Postgres',
      icon: PostgresKubernetesIcon,
      color: 'from-white to-[#666]',
      config: { region: 'us-east-1', size: 'medium', replicas: 2 },
      cost: '$180/month',
    },
    {
      id: 'redis-1',
      name: 'Redis',
      icon: RedisKubernetesIcon,
      color: 'from-[#10b981] to-[#059669]',
      config: { region: 'us-east-1', size: 'small' },
      cost: '$60/month',
    },
    {
      id: 'r2-1',
      name: 'R2 Bucket',
      icon: CloudflareR2BucketIcon,
      color: 'from-[#999] to-[#666]',
      config: { region: 'global' },
      cost: '$15/month',
    },
    {
      id: 'cdn-1',
      name: 'CDN',
      icon: AwsCloudfrontIcon,
      color: 'from-white to-[#666]',
      config: { region: 'global', distribution: 'enabled' },
      cost: '$25/month',
    },
  ];

  const totalCost = deploymentItems.reduce((sum, item) => {
    const cost = parseInt(item.cost.replace(/[^0-9]/g, ''));
    return sum + cost;
  }, 0);

  const handleReview = () => {
    setIsReviewing(true);
    setTimeout(() => {
      setIsReviewing(false);
    }, 2000);
  };

  return (
    <div className="h-full flex flex-col">
      {/* Header */}
      <div className="px-8 py-6 border-b border-gray-200">
        <h2 className="text-2xl font-bold text-gray-900">Deployment Summary</h2>
        <p className="text-gray-600 mt-1">
          Review your infrastructure before deploying
        </p>
      </div>

      <div className="flex-1 overflow-auto">
        <div className="max-w-4xl mx-auto p-8">
          {/* Summary Stats */}
          <div className="grid grid-cols-1 md:grid-cols-4 gap-6 mb-8">
            <motion.div
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: 0.1 }}
              className="bg-[#151515] border border-[#2a2a2a] rounded-xl p-6"
            >
              <div className="flex items-center gap-3 mb-2">
                <Server className="w-5 h-5 text-white" />
                <span className="text-sm font-medium text-white">Resources</span>
              </div>
              <p className="text-2xl font-bold text-gray-900">{deploymentItems.length}</p>
            </motion.div>

            <motion.div
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: 0.2 }}
              className="bg-[#151515] border border-[#2a2a2a] rounded-xl p-6"
            >
              <div className="flex items-center gap-3 mb-2">
                <Clock className="w-5 h-5 text-[#10b981]" />
                <span className="text-sm font-medium text-[#10b981]">Deploy Time</span>
              </div>
              <p className="text-2xl font-bold text-gray-900">~5 min</p>
            </motion.div>

            <motion.div
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: 0.3 }}
              className="bg-[#151515] border border-[#2a2a2a] rounded-xl p-6"
            >
              <div className="flex items-center gap-3 mb-2">
                <DollarSign className="w-5 h-5 text-white" />
                <span className="text-sm font-medium text-white">Est. Cost</span>
              </div>
              <p className="text-2xl font-bold text-gray-900">${totalCost}/mo</p>
            </motion.div>

            <motion.div
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: 0.4 }}
              className="bg-[#151515] border border-[#2a2a2a] rounded-xl p-6"
            >
              <div className="flex items-center gap-3 mb-2">
                <Zap className="w-5 h-5 text-[#a0a0a0]" />
                <span className="text-sm font-medium text-[#a0a0a0]">Environment</span>
              </div>
              <p className="text-2xl font-bold text-gray-900">Production</p>
            </motion.div>
          </div>

          {/* Resources List */}
          <div className="space-y-4 mb-8">
            <h3 className="font-bold text-gray-900">Resources to Deploy</h3>
            {deploymentItems.map((item, index) => (
              <DeployResourceItem key={item.id} item={item} index={index} />
            ))}
          </div>

          {/* Security Check */}
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.6 }}
            className="bg-[#151515] border border-[#2a2a2a] rounded-xl p-6 mb-8"
          >
            <div className="flex items-start gap-4">
              <CheckCircle className="w-6 h-6 text-[#10b981] flex-shrink-0 mt-1" />
              <div>
                <h4 className="font-bold text-gray-900 mb-2">Security Checks Passed</h4>
                <ul className="space-y-1 text-sm text-gray-600">
                  <li>✓ All resources will be deployed in private subnets</li>
                  <li>✓ Encryption at rest enabled for databases</li>
                  <li>✓ Network policies configured for Kubernetes</li>
                  <li>✓ IAM roles follow least-privilege principle</li>
                </ul>
              </div>
            </div>
          </motion.div>

          {/* Warning */}
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.7 }}
            className="bg-[#1a1a1a] border border-[#f59e0b]/30 rounded-xl p-6 mb-8"
          >
            <div className="flex items-start gap-4">
              <AlertTriangle className="w-6 h-6 text-[#f59e0b] flex-shrink-0 mt-1" />
              <div>
                <h4 className="font-bold text-gray-900 mb-1">Demo Notice</h4>
                <p className="text-sm text-gray-600">
                  This is a demonstration. No actual infrastructure will be provisioned.
                  In a real deployment, these resources would be created in your cloud account.
                </p>
              </div>
            </div>
          </motion.div>

          {/* Action Buttons */}
          <div className="flex gap-4">
            <button
              onClick={handleReview}
              disabled={isReviewing}
              className="px-6 py-3 bg-white border border-gray-300 text-gray-700 rounded-lg font-medium hover:bg-gray-50 transition-all duration-200 disabled:opacity-50"
            >
              {isReviewing ? 'Reviewing...' : 'Review Terraform'}
            </button>
            
            <button
              onClick={onDeploy}
              className="flex-1 px-6 py-3 bg-gradient-to-r from-white to-[#666] text-white rounded-lg font-bold hover:from-white hover:to-[#666] transition-all duration-200 shadow-lg hover:shadow-xl transform hover:-translate-y-0.5"
            >
              Deploy Now
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
