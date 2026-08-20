'use client';

import React, { useState } from 'react';
import { motion } from 'framer-motion';
import { Lock, Key, Globe, Shield, Database, CheckCircle, Eye, EyeOff } from 'lucide-react';

interface Secret {
  key: string;
  value: string;
  backend: 'platform' | 'aws';
}

const mockSecrets: Secret[] = [
  { key: 'DATABASE_PASSWORD', value: '••••••••••••', backend: 'aws' },
  { key: 'API_KEY', value: '••••••••••••', backend: 'platform' },
  { key: 'JWT_SECRET', value: '••••••••••••', backend: 'aws' },
];

const mockVariables = [
  { key: 'LOG_LEVEL', value: 'INFO' },
  { key: 'MAX_CONNECTIONS', value: '100' },
  { key: 'CACHE_TTL', value: '3600' },
];

export default function ServiceConfiguration() {
  const [activeTab, setActiveTab] = useState<'secrets' | 'variables' | 'ingress'>('secrets');
  const [showSecrets, setShowSecrets] = useState(false);

  return (
    <div className="h-full flex flex-col">
      {/* Header */}
      <div className="px-8 py-6 border-b border-gray-200">
        <h2 className="text-2xl font-bold text-gray-900">Service Configuration</h2>
        <p className="text-gray-600 mt-1">
          Configure secrets, environment variables, and ingress
        </p>
      </div>

      <div className="flex-1 overflow-auto">
        <div className="max-w-5xl mx-auto p-8">
          {/* Tab navigation */}
          <div className="flex gap-2 mb-8 bg-gray-100 p-1 rounded-lg">
            <button
              onClick={() => setActiveTab('secrets')}
              className={`
                flex-1 px-6 py-3 rounded-lg font-medium transition-all flex items-center justify-center gap-2
                ${activeTab === 'secrets' 
                  ? 'bg-white text-white shadow-md' 
                  : 'text-gray-600 hover:text-gray-900'
                }
              `}
            >
              <Lock className="w-4 h-4" />
              Secrets
            </button>
            <button
              onClick={() => setActiveTab('variables')}
              className={`
                flex-1 px-6 py-3 rounded-lg font-medium transition-all flex items-center justify-center gap-2
                ${activeTab === 'variables' 
                  ? 'bg-white text-white shadow-md' 
                  : 'text-gray-600 hover:text-gray-900'
                }
              `}
            >
              <Key className="w-4 h-4" />
              Variables
            </button>
            <button
              onClick={() => setActiveTab('ingress')}
              className={`
                flex-1 px-6 py-3 rounded-lg font-medium transition-all flex items-center justify-center gap-2
                ${activeTab === 'ingress' 
                  ? 'bg-white text-white shadow-md' 
                  : 'text-gray-600 hover:text-gray-900'
                }
              `}
            >
              <Globe className="w-4 h-4" />
              Ingress & Domain
            </button>
          </div>

          {/* Secrets tab */}
          {activeTab === 'secrets' && (
            <motion.div
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              className="space-y-6"
            >
              {/* Info banner */}
              <div className="bg-gradient-to-r from-[#151515] to-[#151515] border border-[#2a2a2a] rounded-xl p-6">
                <div className="flex items-start gap-3">
                  <Shield className="w-6 h-6 text-white flex-shrink-0 mt-1" />
                  <div>
                    <h3 className="font-bold text-gray-900 mb-2">Secure Secret Management</h3>
                    <p className="text-sm text-gray-600 mb-3">
                      Choose where to store your secrets and how they&apos;re injected into your application.
                    </p>
                    <div className="flex gap-4 text-sm">
                      <div className="flex items-center gap-2 bg-white rounded-lg px-3 py-2 border border-[#2a2a2a]">
                        <Database className="w-4 h-4 text-white" />
                        <span className="text-gray-700">Platform Database</span>
                      </div>
                      <div className="flex items-center gap-2 bg-white rounded-lg px-3 py-2 border border-[#2a2a2a]">
                        <Database className="w-4 h-4 text-[#a0a0a0]" />
                        <span className="text-gray-700">AWS Secrets Manager</span>
                      </div>
                    </div>
                  </div>
                </div>
              </div>

              {/* Secrets list */}
              <div className="bg-white border border-gray-200 rounded-xl overflow-hidden">
                <div className="bg-gray-50 px-6 py-4 border-b border-gray-200 flex items-center justify-between">
                  <h3 className="font-bold text-gray-900">Configured Secrets</h3>
                  <button
                    onClick={() => setShowSecrets(!showSecrets)}
                    className="text-sm text-gray-600 hover:text-gray-900 flex items-center gap-2"
                  >
                    {showSecrets ? (
                      <>
                        <EyeOff className="w-4 h-4" />
                        Hide values
                      </>
                    ) : (
                      <>
                        <Eye className="w-4 h-4" />
                        Show values
                      </>
                    )}
                  </button>
                </div>
                <div className="divide-y divide-gray-200">
                  {mockSecrets.map((secret, index) => (
                    <motion.div
                      key={secret.key}
                      initial={{ opacity: 0, x: -20 }}
                      animate={{ opacity: 1, x: 0 }}
                      transition={{ delay: index * 0.1 }}
                      className="p-6 hover:bg-gray-50 transition-colors"
                    >
                      <div className="flex items-center justify-between">
                        <div className="flex-1">
                          <div className="font-mono text-sm font-semibold text-gray-900 mb-1">
                            {secret.key}
                          </div>
                          <div className="font-mono text-sm text-gray-600">
                            {showSecrets ? 'sk_live_abc123...' : secret.value}
                          </div>
                        </div>
                        <div className="flex items-center gap-3">
                          <span className={`
                            px-3 py-1 rounded-full text-xs font-semibold
                            ${secret.backend === 'aws' 
                              ? 'bg-[#1a1a1a] text-[#a0a0a0]' 
                              : 'bg-[#1a1a1a] text-white'
                            }
                          `}>
                            {secret.backend === 'aws' ? 'AWS Secrets Manager' : 'Platform DB'}
                          </span>
                          <Lock className="w-4 h-4 text-gray-400" />
                        </div>
                      </div>
                    </motion.div>
                  ))}
                </div>
              </div>

              {/* Injection method */}
              <div className="bg-[#151515] border border-[#2a2a2a] rounded-xl p-6">
                <div className="flex items-start gap-3">
                  <CheckCircle className="w-6 h-6 text-[#10b981] flex-shrink-0 mt-1" />
                  <div>
                    <h4 className="font-bold text-gray-900 mb-2">Secrets Injection</h4>
                    <p className="text-sm text-gray-600">
                      Secrets are securely injected as environment variables at deployment time based on your deployment type (ECS, Kubernetes, etc.). 
                      They never appear in logs or version control.
                    </p>
                  </div>
                </div>
              </div>
            </motion.div>
          )}

          {/* Variables tab */}
          {activeTab === 'variables' && (
            <motion.div
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              className="space-y-6"
            >
              <div className="bg-gradient-to-r from-[#151515] to-[#151515] border border-[#2a2a2a] rounded-xl p-6">
                <h3 className="font-bold text-gray-900 mb-2">Environment Variables</h3>
                <p className="text-sm text-gray-600">
                  Non-sensitive configuration values for your application. These are visible in plain text and stored with your service definition.
                </p>
              </div>

              <div className="bg-white border border-gray-200 rounded-xl overflow-hidden">
                <div className="bg-gray-50 px-6 py-4 border-b border-gray-200">
                  <h3 className="font-bold text-gray-900">Configured Variables</h3>
                </div>
                <div className="divide-y divide-gray-200">
                  {mockVariables.map((variable, index) => (
                    <motion.div
                      key={variable.key}
                      initial={{ opacity: 0, x: -20 }}
                      animate={{ opacity: 1, x: 0 }}
                      transition={{ delay: index * 0.1 }}
                      className="p-6 hover:bg-gray-50 transition-colors"
                    >
                      <div className="flex items-center justify-between">
                        <div className="font-mono text-sm font-semibold text-gray-900">
                          {variable.key}
                        </div>
                        <div className="font-mono text-sm text-gray-600 bg-gray-100 px-3 py-1 rounded">
                          {variable.value}
                        </div>
                      </div>
                    </motion.div>
                  ))}
                </div>
              </div>
            </motion.div>
          )}

          {/* Ingress tab */}
          {activeTab === 'ingress' && (
            <motion.div
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              className="space-y-6"
            >
              <div className="bg-gradient-to-r from-[#151515] to-[#151515] border border-[#2a2a2a] rounded-xl p-6">
                <div className="flex items-start gap-3">
                  <Globe className="w-6 h-6 text-white flex-shrink-0 mt-1" />
                  <div>
                    <h3 className="font-bold text-gray-900 mb-2">Ingress Configuration</h3>
                    <p className="text-sm text-gray-600">
                      Platform automatically creates load balancers, listener rules, and SSL certificates based on your ingress configuration.
                    </p>
                  </div>
                </div>
              </div>

              <div className="bg-white border border-gray-200 rounded-xl p-8">
                <h3 className="font-bold text-gray-900 mb-6">Domain Configuration</h3>
                
                <div className="space-y-6">
                  <div>
                    <label className="block text-sm font-semibold text-gray-700 mb-2">
                      Domain Name
                    </label>
                    <div className="flex items-center gap-3">
                      <input
                        type="text"
                        value="billing-api.acmecorp.com"
                        readOnly
                        className="flex-1 px-4 py-3 border border-gray-300 rounded-lg bg-gray-50 font-mono text-sm"
                      />
                      <CheckCircle className="w-5 h-5 text-[#10b981]" />
                    </div>
                  </div>

                  <div>
                    <label className="block text-sm font-semibold text-gray-700 mb-2">
                      Path Pattern
                    </label>
                    <input
                      type="text"
                      value="/api/v1/*"
                      readOnly
                      className="w-full px-4 py-3 border border-gray-300 rounded-lg bg-gray-50 font-mono text-sm"
                    />
                  </div>

                  <div>
                    <label className="block text-sm font-semibold text-gray-700 mb-2">
                      SSL Certificate
                    </label>
                    <div className="flex items-center gap-3 px-4 py-3 border border-gray-300 rounded-lg bg-gray-50">
                      <Lock className="w-5 h-5 text-[#10b981]" />
                      <span className="text-sm text-gray-700">Auto-managed (Let&apos;s Encrypt)</span>
                      <CheckCircle className="w-4 h-4 text-[#10b981] ml-auto" />
                    </div>
                  </div>
                </div>
              </div>

              {/* Auto-created resources */}
              <div className="bg-[#151515] border border-[#2a2a2a] rounded-xl p-6">
                <div className="flex items-start gap-3">
                  <CheckCircle className="w-6 h-6 text-[#10b981] flex-shrink-0 mt-1" />
                  <div>
                    <h4 className="font-bold text-gray-900 mb-2">Automatically Created</h4>
                    <p className="text-sm text-gray-600 mb-3">
                      Platform will automatically create and configure:
                    </p>
                    <ul className="space-y-1 text-sm text-gray-600">
                      <li>✓ Application Load Balancer with health checks</li>
                      <li>✓ Target groups pointing to your service</li>
                      <li>✓ Listener rules for path-based routing</li>
                      <li>✓ SSL/TLS certificate from Let&apos;s Encrypt</li>
                      <li>✓ DNS records in Route 53 or Cloudflare</li>
                    </ul>
                  </div>
                </div>
              </div>
            </motion.div>
          )}
        </div>
      </div>
    </div>
  );
}

