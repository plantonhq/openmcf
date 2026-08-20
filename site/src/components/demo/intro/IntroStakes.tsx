'use client';

import React from 'react';
import { motion } from 'framer-motion';
import {
  XCircle,
  AlertCircle,
  Users,
  Shield,
  Copy,
  Wrench,
  GitBranch,
  Zap,
  Clock,
  TrendingUp,
} from 'lucide-react';

const DevOpsProblems = [
  {
    icon: Wrench,
    title: 'Tool Fragmentation',
    description: '10+ different tools for CI/CD, IaC, monitoring, each with its own learning curve',
  },
  {
    icon: Copy,
    title: 'Copy-Paste Hell',
    description:
      "Same infrastructure code duplicated across environments, teams can't share modules",
  },
  {
    icon: Shield,
    title: 'Security Afterthought',
    description: 'Compliance and governance bolted on later, causing delays and vulnerabilities',
  },
  {
    icon: Users,
    title: 'Tribal Knowledge',
    description: 'Only a few engineers understand the entire system, creating bottlenecks',
  },
];

export default function IntroStakes() {
  return (
    <div className="h-full flex items-center bg-gradient-to-br from-gray-50 to-gray-100 relative rounded-2xl border border-gray-200 shadow-lg overflow-hidden">
      {/* Subtle Background Elements */}
      <div className="absolute inset-0 opacity-20">
        <div className="absolute top-10 left-10 w-32 h-32 bg-red-500/10 rounded-full blur-2xl"></div>
        <div className="absolute bottom-10 right-10 w-24 h-24 bg-white/10 rounded-full blur-2xl"></div>
      </div>

      <div className="max-w-7xl mx-auto px-8 w-full">
        <div className="grid lg:grid-cols-2 gap-4 items-center">
          {/* Left Column - Problems */}
          <motion.div
            initial={{ opacity: 0, x: -20 }}
            animate={{ opacity: 1, x: 0 }}
            transition={{ duration: 0.6 }}
            className="space-y-6"
          >
            {/* Header */}
            <div className="space-y-4">
              <h2 className="text-3xl md:text-4xl font-bold text-gray-900 leading-tight">
                Current DevOps is a{' '}
                <span className="bg-gradient-to-r from-red-600 to-red-500 bg-clip-text text-transparent">
                  nightmare
                </span>
              </h2>
            </div>

            {/* Problems List */}
            <div className="space-y-4">
              {DevOpsProblems.map((problem, index) => (
                <motion.div
                  key={index}
                  initial={{ opacity: 0, x: -20 }}
                  animate={{ opacity: 1, x: 0 }}
                  transition={{ delay: 0.1 + index * 0.1, duration: 0.6 }}
                  className="flex gap-4 p-4 bg-white/60 backdrop-blur-sm border border-gray-200/50 rounded-xl hover:bg-white/80 transition-all duration-300 group"
                >
                  <div className="flex-shrink-0">
                    <div className="w-10 h-10 bg-red-100 border border-red-200 rounded-lg flex items-center justify-center group-hover:bg-red-200 transition-colors duration-300">
                      <XCircle className="w-5 h-5 text-red-600" />
                    </div>
                  </div>

                  <div className="flex-1 space-y-2">
                    <div className="flex items-center gap-2">
                      <problem.icon className="w-5 h-5 text-gray-700" />
                      <h3 className="font-bold text-gray-900 group-hover:text-red-600 transition-colors duration-300">
                        {problem.title}
                      </h3>
                    </div>
                    <p className="text-sm text-gray-600 leading-relaxed">{problem.description}</p>
                  </div>
                </motion.div>
              ))}
            </div>
          </motion.div>

          {/* Right Column - Enhanced Status Dashboard */}
          <motion.div
            initial={{ opacity: 0, scale: 0.9 }}
            animate={{ opacity: 1, scale: 1 }}
            transition={{ duration: 0.6, delay: 0.3 }}
            className="relative"
          >
            <div className="bg-gradient-to-br from-red-50 to-[#151515] border border-red-200 rounded-2xl p-6 shadow-xl">
              <div className="space-y-2">
                {/* Pipeline Visualization */}
                <div className="space-y-5">
                  <div className="flex items-center justify-between">
                    <h3 className="text-xl font-bold text-gray-900">Pipeline Status</h3>
                    <div className="flex items-center gap-2 px-4 py-2 bg-red-600 text-white rounded-lg shadow-md">
                      <XCircle className="w-4 h-4" />
                      <span className="text-sm font-semibold">Pipeline Failed</span>
                    </div>
                  </div>

                  {/* Enhanced Pipeline Flow */}
                  <div className="bg-white/60 backdrop-blur-sm border border-gray-200 rounded-xl p-5">
                    <div className="flex items-center justify-between">
                      {/* Source Stage */}
                      <div className="flex flex-col items-center gap-3">
                        <div className="w-12 h-12 bg-[#1a1a1a] border-2 border-[#2a2a2a] rounded-xl shadow-md flex items-center justify-center">
                          <GitBranch className="w-5 h-5 text-[#10b981]" />
                        </div>
                        <div className="text-center">
                          <span className="text-sm font-semibold text-gray-900">Source</span>
                          <div className="text-xs text-[#10b981] font-medium">✓ Completed</div>
                        </div>
                      </div>

                      {/* Connection Line */}
                      <div className="flex-1 mx-6 relative">
                        <div className="h-1 bg-[#10b981] rounded-full"></div>
                        <div className="absolute top-0 left-0 w-1/2 h-1 bg-[#10b981] rounded-full"></div>
                        <div className="absolute top-0 right-0 w-1/2 h-1 bg-red-300 rounded-full"></div>
                      </div>

                      {/* Build Stage */}
                      <div className="flex flex-col items-center gap-3">
                        <div className="w-12 h-12 bg-red-100 border-2 border-red-300 rounded-xl shadow-md flex items-center justify-center relative">
                          <XCircle className="w-5 h-5 text-red-600" />
                          <div className="absolute -top-1 -right-1 w-4 h-4 bg-red-600 rounded-full flex items-center justify-center">
                            <XCircle className="w-2 h-2 text-white" />
                          </div>
                        </div>
                        <div className="text-center">
                          <span className="text-sm font-semibold text-gray-900">Build</span>
                          <div className="text-xs text-red-600 font-medium">✗ Failed</div>
                        </div>
                      </div>

                      {/* Connection Line */}
                      <div className="flex-1 mx-6 relative">
                        <div className="h-1 bg-red-300 rounded-full"></div>
                        <div className="absolute top-0 left-0 w-full h-1 bg-red-300 rounded-full"></div>
                      </div>

                      {/* Deploy Stage */}
                      <div className="flex flex-col items-center gap-3">
                        <div className="w-12 h-12 bg-gray-100 border-2 border-gray-300 rounded-xl shadow-md flex items-center justify-center opacity-50">
                          <Zap className="w-5 h-5 text-gray-400" />
                        </div>
                        <div className="text-center">
                          <span className="text-sm font-semibold text-gray-400">Deploy</span>
                          <div className="text-xs text-gray-400 font-medium">⏸ Blocked</div>
                        </div>
                      </div>
                    </div>

                    {/* Failure Details */}
                    <div className="mt-4 p-3 bg-red-50 border border-red-200 rounded-lg">
                      <p className="text-sm text-gray-600 leading-relaxed">
                        Environment variables mismatch between dev and staging environments
                      </p>
                    </div>
                  </div>
                </div>

                {/* Team Impact */}
                <div className="space-y-3">
                  <div className="grid grid-cols-1 gap-3">
                    <div className="flex items-center gap-3 p-3 bg-white/60 backdrop-blur-sm border border-red-200 rounded-lg">
                      <div className="w-8 h-8 bg-red-100 rounded-lg flex items-center justify-center">
                        <Users className="w-4 h-4 text-red-600" />
                      </div>
                      <div className="flex-1">
                        <div className="text-sm font-medium text-gray-900">Development Team</div>
                        <div className="text-xs text-red-600">Blocked - 3 features delayed</div>
                      </div>
                      <Clock className="w-4 h-4 text-gray-400" />
                    </div>

                    <div className="flex items-center gap-3 p-3 bg-white/60 backdrop-blur-sm border border-[#2a2a2a] rounded-lg">
                      <div className="w-8 h-8 bg-[#1a1a1a] rounded-lg flex items-center justify-center">
                        <AlertCircle className="w-4 h-4 text-[#a0a0a0]" />
                      </div>
                      <div className="flex-1">
                        <div className="text-sm font-medium text-gray-900">Operations Team</div>
                        <div className="text-xs text-[#a0a0a0]">Firefighting - 2 incidents</div>
                      </div>
                      <TrendingUp className="w-4 h-4 text-gray-400" />
                    </div>

                    <div className="flex items-center gap-3 p-3 bg-white/60 backdrop-blur-sm border border-red-200 rounded-lg">
                      <div className="w-8 h-8 bg-red-100 rounded-lg flex items-center justify-center">
                        <Shield className="w-4 h-4 text-red-600" />
                      </div>
                      <div className="flex-1">
                        <div className="text-sm font-medium text-gray-900">Security Team</div>
                        <div className="text-xs text-red-600">Escalating - 1 vulnerability</div>
                      </div>
                      <AlertCircle className="w-4 h-4 text-gray-400" />
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </motion.div>
        </div>
      </div>
    </div>
  );
}
