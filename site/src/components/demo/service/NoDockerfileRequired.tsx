'use client';

import React, { useState } from 'react';
import { motion } from 'framer-motion';
import { Package, FileCode, CheckCircle, Sparkles, Zap } from 'lucide-react';

const buildMethods = [
  {
    id: 'buildpacks',
    title: 'No Dockerfile Required',
    icon: Zap,
    description: 'Automatic containerization using Cloud-Native BuildPacks',
    benefits: [
      'Zero Dockerfile knowledge needed',
      'Auto-detect tech stack',
      'Best practices built-in',
      'Automatic security updates',
    ],
    status: 'Recommended - Detected Java Spring Boot',
    recommended: true,
  },
  {
    id: 'dockerfile',
    title: 'Bring Your Dockerfile',
    icon: FileCode,
    description: 'Use your existing Dockerfile for full control',
    benefits: [
      'Complete control over build steps',
      'Custom base images',
      'Multi-stage builds',
      'Existing Dockerfile compatibility',
    ],
    status: 'No Dockerfile detected in repository',
  },
];

const supportedLanguages = [
  { name: 'Java', icon: '☕', color: 'from-[#999] to-[#666]' },
  { name: 'Python', icon: '🐍', color: 'from-white to-[#666]' },
  { name: 'Node.js', icon: '🟢', color: 'from-[#10b981] to-[#059669]' },
  { name: 'Go', icon: '🔷', color: 'from-[#999] to-[#666]' },
  { name: 'Ruby', icon: '💎', color: 'from-red-400 to-red-600' },
  { name: 'PHP', icon: '🐘', color: 'from-[#999] to-[#666]' },
  { name: '.NET', icon: '🔵', color: 'from-[#999] to-[#666]' },
  { name: 'Rust', icon: '🦀', color: 'from-white to-[#666]' },
];

export default function NoDockerfileRequired() {
  const [selectedMethod, setSelectedMethod] = useState<string>('buildpacks');

  return (
    <div className="h-full flex flex-col">
      {/* Header */}
      <div className="px-8 py-6 border-b border-gray-200 bg-gradient-to-r from-[#151515] to-[#151515]">
        <motion.div
          initial={{ opacity: 0, y: -20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.6 }}
        >
          <h2 className="text-3xl font-bold text-gray-900 mb-2">No Dockerfile Required</h2>
          <p className="text-gray-600 text-lg">
            Deploy containerized applications without writing or maintaining Dockerfiles
          </p>
        </motion.div>
      </div>

      <div className="flex-1 overflow-auto">
        <div className="max-w-6xl mx-auto p-8">
          {/* Value Prop Banner */}
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.2 }}
            className="bg-gradient-to-r from-white to-[#666] rounded-2xl p-6 text-white mb-8"
          >
            <div className="flex items-center gap-4">
              <div className="w-12 h-12 bg-white/20 rounded-xl flex items-center justify-center backdrop-blur-sm">
                <Zap className="w-6 h-6" />
              </div>
              <div>
                <h3 className="text-xl font-bold mb-1">Like Vercel, But for Backend</h3>
                <p className="text-white/90 text-sm">
                  Push your code and Service Hub handles containerization, deployment, and scaling—all in your cloud.
                </p>
              </div>
            </div>
          </motion.div>

          {/* Build method selection */}
          <div className="grid md:grid-cols-2 gap-6 mb-8">
            {buildMethods.map((method, index) => {
              const Icon = method.icon;
              return (
                <motion.div
                  key={method.id}
                  initial={{ opacity: 0, y: 20 }}
                  animate={{ opacity: 1, y: 0 }}
                  transition={{ delay: 0.3 + index * 0.1 }}
                  onClick={() => setSelectedMethod(method.id)}
                  className={`
                    relative bg-white border-2 rounded-2xl p-8 cursor-pointer transition-all
                    ${selectedMethod === method.id 
                      ? 'border-[#3a3a3a] shadow-xl' 
                      : 'border-gray-200 hover:border-gray-300 hover:shadow-lg'
                    }
                  `}
                >
                  {method.recommended && (
                    <div className="absolute -top-3 right-6 px-3 py-1 bg-gradient-to-r from-white to-[#666] text-white text-xs font-bold rounded-full shadow-lg flex items-center gap-1">
                      <Sparkles className="w-3 h-3" />
                      Recommended
                    </div>
                  )}

                  <div className="flex items-start gap-4 mb-6">
                    <div className={`
                      w-14 h-14 rounded-xl flex items-center justify-center
                      ${selectedMethod === method.id 
                        ? 'bg-gradient-to-br from-white to-[#666]' 
                        : 'bg-gradient-to-br from-gray-400 to-gray-500'
                      }
                    `}>
                      <Icon className="w-7 h-7 text-white" />
                    </div>
                    <div className="flex-1">
                      <h3 className="font-bold text-gray-900 text-xl mb-2">
                        {method.title}
                      </h3>
                      <p className="text-gray-600 text-sm">
                        {method.description}
                      </p>
                    </div>
                    {selectedMethod === method.id && (
                      <motion.div
                        initial={{ scale: 0 }}
                        animate={{ scale: 1 }}
                        className="flex-shrink-0"
                      >
                        <CheckCircle className="w-6 h-6 text-white" />
                      </motion.div>
                    )}
                  </div>

                  <div className="space-y-2 mb-6">
                    {method.benefits.map((benefit, i) => (
                      <div key={i} className="flex items-center gap-2 text-sm text-gray-700">
                        <CheckCircle className="w-4 h-4 text-[#10b981] flex-shrink-0" />
                        <span>{benefit}</span>
                      </div>
                    ))}
                  </div>

                  <div className={`
                    px-4 py-2 rounded-lg text-sm font-medium text-center
                    ${method.recommended 
                      ? 'bg-[#151515] text-[#10b981] border border-[#2a2a2a]' 
                      : 'bg-gray-50 text-gray-600 border border-gray-200'
                    }
                  `}>
                    {method.status}
                  </div>
                </motion.div>
              );
            })}
          </div>

          {/* BuildPacks explanation */}
          {selectedMethod === 'buildpacks' && (
            <motion.div
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              className="space-y-6"
            >
              <div className="bg-gradient-to-r from-[#151515] to-[#151515] border border-[#2a2a2a] rounded-xl p-6">
                <h3 className="font-bold text-gray-900 mb-3 flex items-center gap-2">
                  <Package className="w-5 h-5 text-white" />
                  How It Works: Cloud-Native BuildPacks
                </h3>
                <p className="text-gray-700 mb-4">
                  BuildPacks automatically transform your application source code into production-ready container images—no Dockerfile knowledge required. 
                  They&apos;re open source, maintained by the CNCF community, and follow best practices for security and performance.
                </p>
                <div className="bg-white rounded-lg p-4 border border-[#2a2a2a]">
                  <p className="text-sm text-gray-600">
                    <strong className="text-gray-900">The magic:</strong> Push your code → BuildPacks detect language/framework → 
                    Install dependencies → Compile → Optimize → Create container image. All automatically.
                  </p>
                </div>
              </div>

              {/* Supported languages */}
              <div className="bg-white border border-gray-200 rounded-xl p-8">
                <h3 className="font-bold text-gray-900 mb-6 text-center text-lg">
                  Supported Languages & Frameworks
                </h3>
                <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                  {supportedLanguages.map((lang, index) => (
                    <motion.div
                      key={lang.name}
                      initial={{ opacity: 0, scale: 0.9 }}
                      animate={{ opacity: 1, scale: 1 }}
                      transition={{ delay: 0.5 + index * 0.05 }}
                      className={`
                        bg-gradient-to-br ${lang.color}
                        rounded-xl p-4 text-center text-white shadow-lg
                        hover:scale-105 transition-transform cursor-pointer
                      `}
                    >
                      <div className="text-3xl mb-2">{lang.icon}</div>
                      <div className="font-bold text-sm">{lang.name}</div>
                    </motion.div>
                  ))}
                </div>
                <p className="text-center text-sm text-gray-500 mt-6">
                  Can&apos;t find your tech stack? Create custom BuildPacks for your specific needs.
                </p>
              </div>

              {/* Selected configuration */}
              <motion.div
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: 0.7 }}
                className="bg-[#151515] border border-[#2a2a2a] rounded-xl p-6"
              >
                <div className="flex items-start gap-3">
                  <CheckCircle className="w-6 h-6 text-[#10b981] flex-shrink-0 mt-1" />
                  <div>
                    <h4 className="font-bold text-gray-900 mb-2">Auto-Detected: Java Spring Boot</h4>
                    <p className="text-sm text-gray-600 mb-3">
                      Platform will automatically:
                    </p>
                    <ul className="space-y-1 text-sm text-gray-600">
                      <li>✓ Detect Java 17 runtime</li>
                      <li>✓ Resolve Maven dependencies</li>
                      <li>✓ Apply Spring Boot optimizations</li>
                      <li>✓ Tune JVM for container environment</li>
                      <li>✓ Create optimized, secure container image</li>
                    </ul>
                  </div>
                </div>
              </motion.div>
            </motion.div>
          )}

          {/* Dockerfile option */}
          {selectedMethod === 'dockerfile' && (
            <motion.div
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              className="bg-[#1a1a1a] border border-[#f59e0b]/30 rounded-xl p-6"
            >
              <div className="flex items-start gap-3">
                <FileCode className="w-6 h-6 text-[#f59e0b] flex-shrink-0 mt-1" />
                <div>
                  <h4 className="font-bold text-gray-900 mb-2">Dockerfile Not Found</h4>
                  <p className="text-sm text-gray-600 mb-3">
                    No Dockerfile was detected in your repository. You can either:
                  </p>
                  <ul className="space-y-2 text-sm text-gray-600">
                    <li>• Create a Dockerfile in your repository root</li>
                    <li>• Or switch to &ldquo;No Dockerfile Required&rdquo; for automatic containerization</li>
                  </ul>
                </div>
              </div>
            </motion.div>
          )}
        </div>
      </div>
    </div>
  );
}

