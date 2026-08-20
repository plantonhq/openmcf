'use client';

import React, { useState } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { Edit3, CheckCircle, AlertCircle, ArrowRight } from 'lucide-react';

interface ConfigField {
  label: string;
  current: string;
  modified: string;
  changed: boolean;
}

const initialConfig: ConfigField[] = [
  { label: 'Node Count', current: '2', modified: '2', changed: false },
  { label: 'Instance Type', current: 't3.small', modified: 't3.small', changed: false },
  { label: 'Region', current: 'us-east-1', modified: 'us-east-1', changed: false },
  { label: 'Auto Scaling', current: 'disabled', modified: 'disabled', changed: false },
];

export default function InfraEditFlow() {
  const [config, setConfig] = useState<ConfigField[]>(initialConfig);
  const [isEditing, setIsEditing] = useState(false);
  const [isApplying, setIsApplying] = useState(false);
  const [showSuccess, setShowSuccess] = useState(false);

  const handleEdit = () => {
    setIsEditing(true);
    // Simulate making changes
    setTimeout(() => {
      setConfig([
        { label: 'Node Count', current: '2', modified: '3', changed: true },
        { label: 'Instance Type', current: 't3.small', modified: 't3.medium', changed: true },
        { label: 'Region', current: 'us-east-1', modified: 'us-east-1', changed: false },
        { label: 'Auto Scaling', current: 'disabled', modified: 'enabled', changed: true },
      ]);
    }, 800);
  };

  const handleApply = () => {
    setIsApplying(true);
    setTimeout(() => {
      setIsApplying(false);
      setShowSuccess(true);
      setTimeout(() => {
        setShowSuccess(false);
        setIsEditing(false);
        // Update current values
        setConfig(config.map(field => ({
          ...field,
          current: field.modified,
          changed: false,
        })));
      }, 2000);
    }, 2000);
  };

  const changedCount = config.filter(f => f.changed).length;

  return (
    <div className="h-full flex flex-col">
      {/* Header */}
      <div className="px-8 py-6 border-b border-gray-200">
        <h2 className="text-2xl font-bold text-gray-900">Edit Infrastructure</h2>
        <p className="text-gray-600 mt-1">
          Modify your configuration and redeploy with one click
        </p>
      </div>

      <div className="flex-1 overflow-auto">
        <div className="max-w-6xl mx-auto p-8">
          {/* Info banner */}
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            className="bg-gradient-to-r from-[#151515] to-[#151515] border border-[#2a2a2a] rounded-xl p-6 mb-8"
          >
            <h3 className="font-bold text-gray-900 mb-2">Live Configuration Editing</h3>
            <p className="text-sm text-gray-600">
              Make changes to your infrastructure configuration. The platform detects changes and triggers automated redeployment with full audit trail.
            </p>
          </motion.div>

          {/* Resource card */}
          <div className="bg-white border border-gray-200 rounded-xl shadow-lg overflow-hidden mb-6">
            <div className="bg-gradient-to-r from-white to-[#666] px-6 py-4 flex items-center justify-between">
              <div>
                <h3 className="text-white font-bold text-lg">EKS Cluster Configuration</h3>
                <p className="text-gray-300 text-sm">Production environment</p>
              </div>
              {!isEditing && !showSuccess && (
                <button
                  onClick={handleEdit}
                  className="px-4 py-2 bg-white text-[#a0a0a0] rounded-lg font-medium hover:bg-[#151515] transition-colors flex items-center gap-2"
                >
                  <Edit3 className="w-4 h-4" />
                  Edit Configuration
                </button>
              )}
            </div>

            <div className="p-6">
              {/* Comparison view */}
              <div className="grid grid-cols-2 gap-6">
                {/* Current configuration */}
                <div>
                  <div className="flex items-center gap-2 mb-4">
                    <div className="w-3 h-3 bg-white rounded-full" />
                    <h4 className="font-bold text-gray-900">Current Configuration</h4>
                  </div>
                  <div className="space-y-3">
                    {config.map((field, index) => (
                      <motion.div
                        key={field.label}
                        initial={{ opacity: 0, x: -20 }}
                        animate={{ opacity: 1, x: 0 }}
                        transition={{ delay: index * 0.1 }}
                        className={`
                          bg-gray-50 border rounded-lg p-4
                          ${field.changed ? 'border-red-300 bg-red-50' : 'border-gray-200'}
                        `}
                      >
                        <div className="text-xs font-semibold text-gray-500 mb-1">
                          {field.label}
                        </div>
                        <div className="text-gray-900 font-mono text-sm">
                          {field.current}
                        </div>
                      </motion.div>
                    ))}
                  </div>
                </div>

                {/* Modified configuration */}
                <div>
                  <div className="flex items-center gap-2 mb-4">
                    <div className={`w-3 h-3 rounded-full ${isEditing ? 'bg-[#10b981] animate-pulse' : 'bg-gray-300'}`} />
                    <h4 className="font-bold text-gray-900">
                      {isEditing ? 'Modified Configuration' : 'No Changes'}
                    </h4>
                    {changedCount > 0 && (
                      <span className="ml-auto px-2 py-1 bg-[#1a1a1a] text-[#a0a0a0] text-xs font-semibold rounded-full">
                        {changedCount} changes
                      </span>
                    )}
                  </div>
                  <div className="space-y-3">
                    {config.map((field, index) => (
                      <motion.div
                        key={field.label}
                        initial={{ opacity: 0, x: 20 }}
                        animate={{ opacity: 1, x: 0 }}
                        transition={{ delay: index * 0.1 }}
                        className={`
                          border rounded-lg p-4 relative
                          ${field.changed 
                            ? 'border-[#2a2a2a] bg-[#151515]' 
                            : 'border-gray-200 bg-gray-50'
                          }
                        `}
                      >
                        {field.changed && (
                          <div className="absolute -left-2 top-1/2 transform -translate-y-1/2">
                            <ArrowRight className="w-4 h-4 text-[#10b981]" />
                          </div>
                        )}
                        <div className="text-xs font-semibold text-gray-500 mb-1">
                          {field.label}
                        </div>
                        <div className={`font-mono text-sm ${field.changed ? 'text-[#10b981] font-semibold' : 'text-gray-900'}`}>
                          {field.modified}
                        </div>
                      </motion.div>
                    ))}
                  </div>
                </div>
              </div>

              {/* Change detection alert */}
              <AnimatePresence>
                {isEditing && changedCount > 0 && !isApplying && !showSuccess && (
                  <motion.div
                    initial={{ opacity: 0, y: 20 }}
                    animate={{ opacity: 1, y: 0 }}
                    exit={{ opacity: 0, y: -20 }}
                    className="mt-6 bg-[#1a1a1a] border border-[#f59e0b]/30 rounded-lg p-4 flex items-start gap-3"
                  >
                    <AlertCircle className="w-5 h-5 text-[#f59e0b] flex-shrink-0 mt-0.5" />
                    <div className="flex-1">
                      <h4 className="font-semibold text-gray-900 mb-1">
                        Changes Detected
                      </h4>
                      <p className="text-sm text-gray-600 mb-3">
                        {changedCount} configuration field{changedCount > 1 ? 's have' : ' has'} been modified. Applying these changes will trigger a new deployment.
                      </p>
                      <button
                        onClick={handleApply}
                        className="px-6 py-2 bg-gradient-to-r from-white to-[#666] text-white rounded-lg font-medium hover:from-white hover:to-[#666] transition-colors shadow-lg"
                      >
                        Apply Changes & Redeploy
                      </button>
                    </div>
                  </motion.div>
                )}

                {isApplying && (
                  <motion.div
                    initial={{ opacity: 0, y: 20 }}
                    animate={{ opacity: 1, y: 0 }}
                    className="mt-6 bg-[#151515] border border-[#2a2a2a] rounded-lg p-4 flex items-center gap-3"
                  >
                    <div className="w-5 h-5 border-2 border-[#3a3a3a] border-t-transparent rounded-full animate-spin" />
                    <span className="text-sm font-medium text-gray-900">
                      Applying changes and redeploying infrastructure...
                    </span>
                  </motion.div>
                )}

                {showSuccess && (
                  <motion.div
                    initial={{ opacity: 0, scale: 0.9 }}
                    animate={{ opacity: 1, scale: 1 }}
                    className="mt-6 bg-[#151515] border border-[#2a2a2a] rounded-lg p-4 flex items-start gap-3"
                  >
                    <CheckCircle className="w-5 h-5 text-[#10b981] flex-shrink-0 mt-0.5" />
                    <div>
                      <h4 className="font-semibold text-gray-900 mb-1">
                        Changes Applied Successfully
                      </h4>
                      <p className="text-sm text-gray-600">
                        Your infrastructure has been updated. All changes are now live and tracked in version history.
                      </p>
                    </div>
                  </motion.div>
                )}
              </AnimatePresence>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

