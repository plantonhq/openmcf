'use client';

import React from 'react';
import { motion } from 'framer-motion';
import { Clock, User, GitCommit, RotateCcw, ChevronRight } from 'lucide-react';

interface VersionHistoryItem {
  id: string;
  timestamp: string;
  user: string;
  action: string;
  changes: {
    resource: string;
    field: string;
    oldValue: string;
    newValue: string;
  }[];
}

const versionHistory: VersionHistoryItem[] = [
  {
    id: 'v4',
    timestamp: '2025-10-22 14:30:00',
    user: 'sarah.chen@acmecorp.com',
    action: 'Updated EKS cluster configuration',
    changes: [
      {
        resource: 'EKS Cluster',
        field: 'node_count',
        oldValue: '2',
        newValue: '3',
      },
      {
        resource: 'EKS Cluster',
        field: 'instance_type',
        oldValue: 't3.small',
        newValue: 't3.medium',
      },
    ],
  },
  {
    id: 'v3',
    timestamp: '2025-10-22 10:15:00',
    user: 'john.doe@acmecorp.com',
    action: 'Added Redis instance',
    changes: [
      {
        resource: 'Redis',
        field: 'size',
        oldValue: '-',
        newValue: 'small',
      },
      {
        resource: 'Redis',
        field: 'region',
        oldValue: '-',
        newValue: 'us-east-1',
      },
    ],
  },
  {
    id: 'v2',
    timestamp: '2025-10-21 16:45:00',
    user: 'jane.smith@acmecorp.com',
    action: 'Updated Postgres configuration',
    changes: [
      {
        resource: 'Postgres',
        field: 'replicas',
        oldValue: '1',
        newValue: '2',
      },
    ],
  },
  {
    id: 'v1',
    timestamp: '2025-10-21 14:00:00',
    user: 'john.doe@acmecorp.com',
    action: 'Initial deployment',
    changes: [
      {
        resource: 'VPC',
        field: 'cidr',
        oldValue: '-',
        newValue: '10.0.0.0/16',
      },
      {
        resource: 'EKS Cluster',
        field: 'node_count',
        oldValue: '-',
        newValue: '2',
      },
    ],
  },
];

export default function InfraVersionHistory() {
  return (
    <div className="h-full flex flex-col">
      {/* Header */}
      <div className="px-8 py-6 border-b border-gray-200">
        <h2 className="text-2xl font-bold text-gray-900">Version History</h2>
        <p className="text-gray-600 mt-1">
          Track all changes with complete audit trail
        </p>
      </div>

      <div className="flex-1 overflow-auto">
        <div className="max-w-5xl mx-auto p-8">
          {/* Auditing Benefits */}
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            className="bg-gradient-to-r from-[#151515] to-[#151515] border border-[#2a2a2a] rounded-xl p-6 mb-8"
          >
            <h3 className="font-bold text-gray-900 mb-2">Complete Audit Trail</h3>
            <p className="text-sm text-gray-600">
              Track <strong>who</strong> made <strong>what</strong> change, <strong>when</strong>. Every modification is recorded with git-diff style visualization for easy comparison.
            </p>
          </motion.div>

          {/* Version Timeline */}
          <div className="space-y-6">
            {versionHistory.map((version, index) => (
              <motion.div
                key={version.id}
                initial={{ opacity: 0, x: -20 }}
                animate={{ opacity: 1, x: 0 }}
                transition={{ delay: index * 0.1 }}
                className="relative"
              >
                {/* Timeline connector */}
                {index < versionHistory.length - 1 && (
                  <div className="absolute left-6 top-16 bottom-0 w-0.5 bg-gray-200" />
                )}

                <div className="bg-white border border-gray-200 rounded-xl p-6 hover:shadow-lg transition-shadow">
                  <div className="flex items-start gap-6">
                    {/* Version indicator */}
                    <div className="flex-shrink-0">
                      <div className="w-12 h-12 rounded-full bg-[#1a1a1a] border border-[#2a2a2a] flex items-center justify-center text-white font-bold shadow-lg">
                        {version.id.toUpperCase()}
                      </div>
                    </div>

                    {/* Version details */}
                    <div className="flex-1">
                      <div className="flex items-start justify-between mb-3">
                        <div>
                          <h4 className="font-bold text-gray-900 text-lg mb-1">
                            {version.action}
                          </h4>
                          <div className="flex flex-wrap items-center gap-4 text-sm text-gray-600">
                            <div className="flex items-center gap-1.5">
                              <User className="w-4 h-4" />
                              <span>{version.user}</span>
                            </div>
                            <div className="flex items-center gap-1.5">
                              <Clock className="w-4 h-4" />
                              <span>{version.timestamp}</span>
                            </div>
                          </div>
                        </div>

                        {index === 0 && (
                          <span className="px-3 py-1 bg-[#1a1a1a] text-[#10b981] text-xs font-semibold rounded-full">
                            Current
                          </span>
                        )}
                      </div>

                      {/* Git-diff style changes */}
                      <div className="bg-gray-50 rounded-lg p-4 space-y-3">
                        <div className="flex items-center gap-2 text-sm font-medium text-gray-700 mb-2">
                          <GitCommit className="w-4 h-4" />
                          <span>Changes ({version.changes.length})</span>
                        </div>
                        
                        {version.changes.map((change, changeIndex) => (
                          <div key={changeIndex} className="space-y-1">
                            <div className="text-xs font-semibold text-gray-600 mb-1">
                              {change.resource} → {change.field}
                            </div>
                            <div className="flex items-center gap-2 text-sm font-mono">
                              {change.oldValue !== '-' && (
                                <div className="flex-1 bg-red-50 border border-red-200 rounded px-3 py-1.5 text-red-700">
                                  <span className="text-[#ef4444] mr-2">-</span>
                                  {change.oldValue}
                                </div>
                              )}
                              {change.oldValue !== '-' && <ChevronRight className="w-4 h-4 text-gray-400 flex-shrink-0" />}
                              <div className="flex-1 bg-[#151515] border border-[#2a2a2a] rounded px-3 py-1.5 text-[#10b981]">
                                <span className="text-[#10b981] mr-2">+</span>
                                {change.newValue}
                              </div>
                            </div>
                          </div>
                        ))}
                      </div>

                      {/* Action buttons */}
                      {index > 0 && (
                        <div className="mt-4 flex gap-3">
                          <button className="px-4 py-2 bg-white border border-gray-300 text-gray-700 rounded-lg text-sm font-medium hover:bg-gray-50 transition-colors flex items-center gap-2">
                            <RotateCcw className="w-4 h-4" />
                            Rollback to this version
                          </button>
                          <button className="px-4 py-2 bg-white border border-gray-300 text-gray-700 rounded-lg text-sm font-medium hover:bg-gray-50 transition-colors">
                            View Full Diff
                          </button>
                        </div>
                      )}
                    </div>
                  </div>
                </div>
              </motion.div>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
}

