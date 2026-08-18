'use client';

import React, { useState } from 'react';
import { motion } from 'framer-motion';
import { Building2, Users, Rocket, Check } from 'lucide-react';

export type CompanyType = 'it-consulting' | 'small-product' | 'established-product';

interface CompanyTypeSelectionProps {
  onCompanyTypeSelect: (type: CompanyType) => void;
}

const companyTypes = [
  {
    id: 'it-consulting' as CompanyType,
    name: 'IT Consulting Agency',
    icon: Rocket,
    description: 'Bootstrap infrastructure fast for client projects',
    painPoints: ['3+ weeks per project setup', 'Lost deals due to slow delivery', 'Client onboarding delays'],
    color: 'from-white/10 to-white/5',
  },
  {
    id: 'small-product' as CompanyType,
    name: 'Small Product Company',
    icon: Users,
    description: 'Self-service DevOps for developer teams',
    painPoints: ['Can\'t afford DevOps team', 'Developers stuck in ops work', 'Scaling bottleneck'],
    color: 'from-white/10 to-white/5',
  },
  {
    id: 'established-product' as CompanyType,
    name: 'Established Product Company',
    icon: Building2,
    description: 'Governance, visibility, and compliance at scale',
    painPoints: ['Infrastructure sprawl', 'Compliance gaps', 'Security misconfigurations'],
    color: 'from-white/10 to-white/5',
  },
];

export default function CompanyTypeSelection({ onCompanyTypeSelect }: CompanyTypeSelectionProps) {
  const [selectedType, setSelectedType] = useState<CompanyType | null>(null);

  const handleTypeClick = (type: CompanyType) => {
    setSelectedType(type);
    // Automatically advance to next slide after a brief visual feedback
    setTimeout(() => {
      onCompanyTypeSelect(type);
    }, 400);
  };

  return (
    <div className="min-h-screen bg-[#0a0a0a] flex items-center justify-center p-8">
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.6 }}
        className="text-center max-w-6xl mx-auto"
      >
        <motion.div
          initial={{ scale: 0 }}
          animate={{ scale: 1 }}
          transition={{ duration: 0.5 }}
          className="inline-flex items-center gap-2 bg-gradient-to-r from-[#111] to-[#222] text-white px-6 py-3 rounded-full text-lg font-bold mb-8"
        >
          <Building2 className="w-5 h-5" />
          Personalized Demo
        </motion.div>

        <h1 className="text-4xl md:text-5xl lg:text-6xl font-black text-white mb-6 tracking-tight">
          What type of company are you?
        </h1>

        <p className="text-xl text-gray-300 mb-12 max-w-3xl mx-auto">
          We&apos;ll customize this demo to show you the features most relevant to your needs.
        </p>

        <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-12">
          {companyTypes.map((type, index) => {
            const Icon = type.icon;
            return (
              <motion.div
                key={type.id}
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: index * 0.1 + 0.2 }}
                onClick={() => handleTypeClick(type.id)}
                className={`relative cursor-pointer transition-all duration-300 ${
                  selectedType === type.id ? 'scale-105' : 'hover:scale-105'
                }`}
              >
                <div
                  className={`bg-[#1a1a1a] border border-[#2a2a2a] rounded-2xl p-8 hover:border-[#3a3a3a] hover:bg-[#1f1f1f] transition-all duration-200 h-full ${
                    selectedType === type.id
                      ? 'ring-2 ring-white/50'
                      : ''
                  }`}
                >
                  {selectedType === type.id && (
                    <motion.div
                      initial={{ scale: 0 }}
                      animate={{ scale: 1 }}
                      className="absolute -top-3 -right-3 w-8 h-8 bg-[#fff] rounded-full flex items-center justify-center"
                    >
                      <Check className="w-5 h-5 text-black" />
                    </motion.div>
                  )}

                  <div className="flex flex-col items-center gap-4 text-center">
                    <div className="w-16 h-16 rounded-2xl bg-white/10 flex items-center justify-center">
                      <Icon className="w-8 h-8 text-white" />
                    </div>
                    <h3 className="font-bold text-xl text-white">{type.name}</h3>
                    <p className="text-[#a0a0a0] text-sm leading-relaxed">{type.description}</p>
                    
                    <div className="mt-4 space-y-2 w-full">
                      {type.painPoints.map((point, i) => (
                        <div key={i} className="text-xs text-[#666] bg-white/5 rounded-lg px-3 py-2">
                          {point}
                        </div>
                      ))}
                    </div>
                  </div>
                </div>
              </motion.div>
            );
          })}
        </div>

        <motion.div
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ delay: 0.6 }}
          className="flex flex-col items-center gap-4"
        >
          <p className="text-sm text-gray-400">
            {selectedType 
              ? 'Starting personalized demo...' 
              : 'Select an option above to continue'}
          </p>
          <p className="text-xs text-gray-500">
            You can switch between different company types anytime by clicking Home
          </p>
        </motion.div>
      </motion.div>
    </div>
  );
}

