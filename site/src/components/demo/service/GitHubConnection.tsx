'use client';

import React, { useState } from 'react';
import { motion } from 'framer-motion';
import { Github, Search, Star, GitBranch, CheckCircle } from 'lucide-react';

interface Repository {
  id: string;
  name: string;
  description: string;
  language: string;
  stars: number;
  isPrivate: boolean;
}

const repositories: Repository[] = [
  {
    id: '1',
    name: 'billing-service',
    description: 'Java Spring Boot billing and payment processing service',
    language: 'Java',
    stars: 12,
    isPrivate: true,
  },
  {
    id: '2',
    name: 'user-api',
    description: 'Python FastAPI user management and authentication',
    language: 'Python',
    stars: 8,
    isPrivate: true,
  },
  {
    id: '3',
    name: 'notification-service',
    description: 'Node.js Express service for email and SMS notifications',
    language: 'TypeScript',
    stars: 15,
    isPrivate: true,
  },
  {
    id: '4',
    name: 'analytics-worker',
    description: 'Go-based background worker for analytics processing',
    language: 'Go',
    stars: 6,
    isPrivate: false,
  },
];

const languageColors: Record<string, string> = {
  Java: 'bg-white',
  Python: 'bg-white',
  TypeScript: 'bg-[#111]',
  Go: 'bg-white',
};

export default function GitHubConnection() {
  const [selectedRepo, setSelectedRepo] = useState<string | null>(null);
  const [isConnecting, setIsConnecting] = useState(false);
  const [connected, setConnected] = useState(false);

  const handleConnect = (repoId: string) => {
    setSelectedRepo(repoId);
    setIsConnecting(true);
    setTimeout(() => {
      setIsConnecting(false);
      setConnected(true);
    }, 1500);
  };

  return (
    <div className="h-full flex flex-col">
      {/* Header */}
      <div className="px-8 py-6 border-b border-gray-200">
        <h2 className="text-2xl font-bold text-gray-900">Connect GitHub Repository</h2>
        <p className="text-gray-600 mt-1">
          Select a repository to deploy
        </p>
      </div>

      <div className="flex-1 overflow-auto">
        <div className="max-w-5xl mx-auto p-8">
          {/* GitHub account card */}
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            className="bg-gradient-to-r from-gray-800 to-gray-900 rounded-xl p-6 mb-8 text-white shadow-lg"
          >
            <div className="flex items-center gap-4">
              <div className="w-16 h-16 bg-white rounded-full flex items-center justify-center">
                <Github className="w-10 h-10 text-gray-900" />
              </div>
              <div className="flex-1">
                <h3 className="font-bold text-lg">ACME Corp Engineering</h3>
                <p className="text-gray-400 text-sm">@acmecorp-engineering</p>
              </div>
              <div className="flex items-center gap-2 text-sm">
                <CheckCircle className="w-5 h-5 text-[#10b981]" />
                <span className="text-[#10b981] font-medium">Connected</span>
              </div>
            </div>
          </motion.div>

          {/* Search bar */}
          <div className="mb-6">
            <div className="relative">
              <Search className="absolute left-4 top-1/2 transform -translate-y-1/2 w-5 h-5 text-gray-400" />
              <input
                type="text"
                placeholder="Search repositories..."
                className="w-full pl-12 pr-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-white focus:border-transparent outline-none"
              />
            </div>
          </div>

          {/* Repositories list */}
          <div className="space-y-4">
            <h3 className="font-bold text-gray-900 flex items-center gap-2">
              <GitBranch className="w-5 h-5" />
              Select Repository
            </h3>

            {repositories.map((repo, index) => (
              <motion.div
                key={repo.id}
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: index * 0.1 }}
                className={`
                  bg-white border-2 rounded-xl p-6 transition-all cursor-pointer
                  ${selectedRepo === repo.id 
                    ? 'border-[#3a3a3a] shadow-lg' 
                    : 'border-gray-200 hover:border-gray-300 hover:shadow-md'
                  }
                `}
                onClick={() => !connected && setSelectedRepo(repo.id)}
              >
                <div className="flex items-start justify-between">
                  <div className="flex-1">
                    <div className="flex items-center gap-3 mb-2">
                      <h4 className="font-bold text-gray-900 text-lg">
                        {repo.name}
                      </h4>
                      <span className={`px-2 py-1 ${repo.isPrivate ? 'bg-gray-100 text-gray-600' : 'bg-[#1a1a1a] text-[#10b981]'} text-xs font-semibold rounded`}>
                        {repo.isPrivate ? 'Private' : 'Public'}
                      </span>
                    </div>
                    <p className="text-gray-600 text-sm mb-3">
                      {repo.description}
                    </p>
                    <div className="flex items-center gap-4 text-sm">
                      <div className="flex items-center gap-1.5">
                        <div className={`w-3 h-3 ${languageColors[repo.language]} rounded-full`} />
                        <span className="text-gray-700">{repo.language}</span>
                      </div>
                      <div className="flex items-center gap-1.5 text-gray-600">
                        <Star className="w-4 h-4" />
                        <span>{repo.stars}</span>
                      </div>
                    </div>
                  </div>

                  <div className="ml-4">
                    {connected && selectedRepo === repo.id ? (
                      <motion.div
                        initial={{ scale: 0 }}
                        animate={{ scale: 1 }}
                        className="w-12 h-12 bg-[#1a1a1a] rounded-full flex items-center justify-center"
                      >
                        <CheckCircle className="w-6 h-6 text-[#10b981]" />
                      </motion.div>
                    ) : selectedRepo === repo.id && !connected ? (
                      <button
                        onClick={() => handleConnect(repo.id)}
                        disabled={isConnecting}
                        className="px-6 py-2 bg-[#fff] text-black rounded-lg font-medium hover:bg-gray-200 transition-colors disabled:opacity-50"
                      >
                        {isConnecting ? (
                          <span className="flex items-center gap-2">
                            <div className="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin" />
                            Connecting...
                          </span>
                        ) : (
                          'Connect Repository'
                        )}
                      </button>
                    ) : null}
                  </div>
                </div>
              </motion.div>
            ))}
          </div>

          {/* Auto-detection info */}
          {connected && (
            <motion.div
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              className="mt-8 bg-[#151515] border border-[#2a2a2a] rounded-xl p-6"
            >
              <div className="flex items-start gap-3">
                <CheckCircle className="w-6 h-6 text-[#10b981] flex-shrink-0 mt-1" />
                <div>
                  <h4 className="font-bold text-gray-900 mb-2">Repository Connected Successfully</h4>
                  <p className="text-sm text-gray-600 mb-3">
                    Platform has analyzed your repository and detected:
                  </p>
                  <ul className="space-y-1 text-sm text-gray-600">
                    <li>✓ Tech Stack: Java 17 with Spring Boot 3.2</li>
                    <li>✓ Build Tool: Maven (pom.xml detected)</li>
                    <li>✓ Dockerfile: Not found (BuildPacks recommended)</li>
                    <li>✓ Dependencies: 24 packages analyzed</li>
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

