'use client';

import React, { useState } from 'react';
import { motion } from 'framer-motion';
import { Network } from 'lucide-react';

interface ResourceNode {
  id: string;
  name: string;
  type: string;
  x: number;
  y: number;
  color: string;
  connections: string[];
}

const resources: ResourceNode[] = [
  {
    id: 'vpc',
    name: 'VPC',
    type: 'Network',
    x: 50,
    y: 20,
    color: 'from-white to-[#666]',
    connections: ['eks', 'postgres', 'redis'],
  },
  {
    id: 'eks',
    name: 'EKS Cluster',
    type: 'Compute',
    x: 25,
    y: 50,
    color: 'from-white to-[#666]',
    connections: ['postgres', 'redis', 'alb'],
  },
  {
    id: 'postgres',
    name: 'Postgres',
    type: 'Database',
    x: 50,
    y: 70,
    color: 'from-white to-[#666]',
    connections: [],
  },
  {
    id: 'redis',
    name: 'Redis',
    type: 'Cache',
    x: 75,
    y: 70,
    color: 'from-[#10b981] to-[#059669]',
    connections: [],
  },
  {
    id: 'alb',
    name: 'Load Balancer',
    type: 'Network',
    x: 25,
    y: 85,
    color: 'from-white to-[#666]',
    connections: ['cdn'],
  },
  {
    id: 'cdn',
    name: 'CloudFront CDN',
    type: 'Content Delivery',
    x: 75,
    y: 50,
    color: 'from-[#999] to-[#666]',
    connections: ['r2'],
  },
  {
    id: 'r2',
    name: 'R2 Bucket',
    type: 'Storage',
    x: 85,
    y: 30,
    color: 'from-[#999] to-[#666]',
    connections: [],
  },
];

export default function InfraVisualization() {
  const [hoveredNode, setHoveredNode] = useState<string | null>(null);
  const [selectedNode, setSelectedNode] = useState<string | null>(null);

  const getNodePosition = (node: ResourceNode) => ({
    left: `${node.x}%`,
    top: `${node.y}%`,
  });

  const selected = selectedNode ? resources.find(r => r.id === selectedNode) : null;

  return (
    <div className="h-full flex flex-col">
      {/* Header */}
      <div className="px-8 py-6 border-b border-gray-200">
        <h2 className="text-2xl font-bold text-gray-900">Infrastructure Visualization</h2>
        <p className="text-gray-600 mt-1">
          See your entire infrastructure and how resources connect
        </p>
      </div>

      <div className="flex-1 overflow-auto bg-gradient-to-br from-gray-50 to-gray-100 relative">
        <div className="max-w-7xl mx-auto p-8 h-full">
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            className="bg-white rounded-2xl shadow-xl p-8 h-full border border-gray-200"
          >
            {/* Info banner */}
            <div className="mb-6 bg-gradient-to-r from-[#151515] to-[#151515] border border-[#2a2a2a] rounded-lg p-4 flex items-start gap-3">
              <Network className="w-5 h-5 text-white flex-shrink-0 mt-0.5" />
              <div>
                <h3 className="font-semibold text-gray-900 mb-1">Interactive Infrastructure Map</h3>
                <p className="text-sm text-gray-600">
                  Click on any resource to see its details and connections. Lines show dependencies between resources.
                </p>
              </div>
            </div>

            {/* Visualization container */}
            <div className="relative h-[calc(100%-8rem)] bg-gradient-to-br from-slate-50 to-slate-100 rounded-xl border-2 border-dashed border-gray-300 overflow-hidden">
              {/* Connection lines (SVG) */}
              <svg className="absolute inset-0 w-full h-full pointer-events-none" style={{ zIndex: 1 }}>
                <defs>
                  <marker
                    id="arrowhead"
                    markerWidth="10"
                    markerHeight="10"
                    refX="9"
                    refY="3"
                    orient="auto"
                  >
                    <polygon points="0 0, 10 3, 0 6" fill="#9ca3af" />
                  </marker>
                </defs>
                {resources.map((from) =>
                  from.connections.map((toId) => {
                    const to = resources.find((r) => r.id === toId);
                    if (!to) return null;
                    
                    const isHighlighted = 
                      hoveredNode === from.id || 
                      hoveredNode === to.id ||
                      selectedNode === from.id ||
                      selectedNode === to.id;

                    return (
                      <line
                        key={`${from.id}-${toId}`}
                        x1={`${from.x}%`}
                        y1={`${from.y}%`}
                        x2={`${to.x}%`}
                        y2={`${to.y}%`}
                        stroke={isHighlighted ? '#ededed' : '#9ca3af'}
                        strokeWidth={isHighlighted ? '3' : '2'}
                        strokeDasharray="5,5"
                        markerEnd="url(#arrowhead)"
                        opacity={isHighlighted ? 1 : 0.4}
                        className="transition-all duration-300"
                      />
                    );
                  })
                )}
              </svg>

              {/* Resource nodes */}
              {resources.map((node, index) => (
                <motion.div
                  key={node.id}
                  initial={{ opacity: 0, scale: 0 }}
                  animate={{ opacity: 1, scale: 1 }}
                  transition={{ delay: index * 0.1 }}
                  style={{
                    ...getNodePosition(node),
                    zIndex: hoveredNode === node.id || selectedNode === node.id ? 10 : 5
                  }}
                  className="absolute transform -translate-x-1/2 -translate-y-1/2 cursor-pointer"
                  onMouseEnter={() => setHoveredNode(node.id)}
                  onMouseLeave={() => setHoveredNode(null)}
                  onClick={() => setSelectedNode(node.id === selectedNode ? null : node.id)}
                >
                  <div
                    className={`
                      relative group
                      ${hoveredNode === node.id || selectedNode === node.id ? 'scale-110' : 'scale-100'}
                      transition-transform duration-200
                    `}
                  >
                    {/* Node circle */}
                    <div
                      className={`
                        w-24 h-24 rounded-full 
                        bg-gradient-to-br ${node.color}
                        shadow-lg
                        flex items-center justify-center
                        border-4 border-white
                        ${selectedNode === node.id ? 'ring-4 ring-white' : ''}
                      `}
                    >
                      <div className="text-center">
                        <div className="text-white font-bold text-sm leading-tight">
                          {node.name}
                        </div>
                        <div className="text-white/80 text-xs mt-1">
                          {node.type}
                        </div>
                      </div>
                    </div>

                    {/* Hover tooltip */}
                    {hoveredNode === node.id && selectedNode !== node.id && (
                      <motion.div
                        initial={{ opacity: 0, y: 10 }}
                        animate={{ opacity: 1, y: 0 }}
                        className="absolute top-full left-1/2 transform -translate-x-1/2 mt-2 bg-gray-900 text-white text-xs rounded-lg px-3 py-2 whitespace-nowrap shadow-xl z-20"
                      >
                        Click to see details
                        <div className="absolute bottom-full left-1/2 transform -translate-x-1/2 w-0 h-0 border-l-4 border-r-4 border-b-4 border-transparent border-b-gray-900" />
                      </motion.div>
                    )}
                  </div>
                </motion.div>
              ))}
            </div>

            {/* Selected node details */}
            {selected && (
              <motion.div
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                className="absolute bottom-8 right-8 bg-white border-2 border-[#3a3a3a] rounded-xl shadow-2xl p-6 w-80"
              >
                <div className="flex items-start justify-between mb-4">
                  <div>
                    <h3 className="font-bold text-gray-900 text-lg">{selected.name}</h3>
                    <p className="text-sm text-gray-600">{selected.type}</p>
                  </div>
                  <button
                    onClick={() => setSelectedNode(null)}
                    className="text-gray-400 hover:text-gray-600"
                  >
                    ×
                  </button>
                </div>
                
                <div className="space-y-3">
                  <div>
                    <div className="text-xs font-semibold text-gray-500 mb-1">STATUS</div>
                    <div className="flex items-center gap-2">
                      <div className="w-2 h-2 bg-[#10b981] rounded-full animate-pulse" />
                      <span className="text-sm text-gray-900">Active</span>
                    </div>
                  </div>
                  
                  <div>
                    <div className="text-xs font-semibold text-gray-500 mb-1">CONNECTIONS</div>
                    <div className="text-sm text-gray-900">
                      {selected.connections.length > 0 
                        ? selected.connections.map(id => 
                            resources.find(r => r.id === id)?.name
                          ).join(', ')
                        : 'No outgoing connections'
                      }
                    </div>
                  </div>

                  <div className="pt-3 border-t border-gray-200">
                    <button className="w-full px-4 py-2 bg-gradient-to-r from-white to-[#666] text-white rounded-lg text-sm font-medium hover:from-white hover:to-[#666] transition-colors">
                      View Configuration
                    </button>
                  </div>
                </div>
              </motion.div>
            )}
          </motion.div>
        </div>
      </div>
    </div>
  );
}

