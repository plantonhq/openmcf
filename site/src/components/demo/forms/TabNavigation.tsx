'use client';

import React from 'react';

export interface TabItem {
  id: string;
  label: string;
  count?: number;
  content?: React.ReactNode;
  icon?: React.ReactNode;
}

interface TabNavigationProps {
  tabs: TabItem[];
  activeTab: string;
  onTabChange: (tabId: string) => void;
  className?: string;
}

export const TabNavigation: React.FC<TabNavigationProps> = ({
  tabs,
  activeTab,
  onTabChange,
  className = '',
}) => {
  return (
    <div className={`w-full ${className}`}>
      {/* Tab Navigation */}
      <div className="flex gap-6" style={{ minHeight: '40px', alignItems: 'center' }}>
        {tabs.map((tab) => {
          const isActive = activeTab === tab.id;
          return (
            <button
              key={tab.id}
              type="button"
              onClick={(e) => {
                e.preventDefault();
                e.stopPropagation();
                onTabChange(tab.id);
              }}
              className="relative transition-all duration-200"
              style={{
                color: '#0E131E',
                textTransform: 'none',
                fontSize: '12px',
                fontWeight: 500,
                padding: '8px 12px',
                borderRadius: '6px',
                backgroundColor: isActive ? '#E7EAED' : 'transparent',
                border: 'none',
                cursor: 'pointer',
                minHeight: 'auto'
              }}
              aria-current={isActive ? 'page' : undefined}
            >
              {isActive && (
                <div 
                  className="absolute left-0 right-0 h-0.5" 
                  style={{ backgroundColor: '#005DA0', bottom: '-12px' }}
                ></div>
              )}
              <span className="flex items-center gap-2">
                {tab.icon && (
                  <span className={isActive ? 'text-gray-900' : 'text-gray-400'}>
                    {tab.icon}
                  </span>
                )}
                <span>
                  {tab.label}
                  {tab.count !== undefined && ` (${tab.count})`}
                </span>
              </span>
            </button>
          );
        })}
      </div>

      {/* Tab Content */}
      {tabs.find(tab => tab.id === activeTab)?.content && (
        <div className="mt-6">
          {tabs.find(tab => tab.id === activeTab)?.content}
        </div>
      )}
    </div>
  );
};

export default TabNavigation;
