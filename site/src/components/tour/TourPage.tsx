'use client';

import React, { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import Image from 'next/image';
import { Button } from '@/components/tour/ui/button';
import { ChevronLeft, ChevronRight, X, Play } from 'lucide-react';
import { motion, AnimatePresence } from 'framer-motion';
import { Toaster } from '@/components/tour/ui/toaster';
import CalloutCard from '@/components/tour/CalloutCard';

type Callout = { number: number; x: number; y: number; title: string; description: string };
type TourStep = {
  step_number: number;
  title: string;
  description: string;
  image_url: string;
  callouts?: Callout[];
};

const TOUR_STEPS: TourStep[] = [
  {
    step_number: 1,
    title: 'Main Dashboard',
    description:
      "Get a bird's-eye view of your entire infrastructure. Monitor resources, recent activities, and key metrics all in one place.",
    image_url:
      'https://qtrypzzcjebvfcihiynt.supabase.co/storage/v1/object/public/base44-prod/public/a30b2df17_Xnapper-2025-08-29-071655.png',
    callouts: [
      {
        number: 1,
        x: 10,
        y: 40,
        title: 'Navigation Menu',
        description:
          'Navigate easily between different hubs like Service Hub, Agent Store, and Connections.',
      },
      {
        number: 2,
        x: 45,
        y: 35,
        title: 'Summary Cards',
        description:
          'Track key metrics such as active users, automation minutes, and total resources.',
      },
      {
        number: 3,
        x: 55,
        y: 80,
        title: 'Recent Activities',
        description:
          'See the latest pipelines and stack jobs at a glance to stay on top of deployments.',
      },
    ],
  },
  {
    step_number: 2,
    title: 'Deployment Component Store',
    description:
      'Browse a rich marketplace of pre-built components to quickly deploy services and infrastructure across any cloud.',
    image_url:
      'https://qtrypzzcjebvfcihiynt.supabase.co/storage/v1/object/public/base44-prod/public/9c682fceb_Xnapper-2025-08-29-071702.png',
    callouts: [
      {
        number: 1,
        x: 15,
        y: 50,
        title: 'Smart Filters',
        description:
          'Filter components by cloud provider or type to find exactly what you need in seconds.',
      },
      {
        number: 2,
        x: 55,
        y: 50,
        title: 'Component Library',
        description:
          'Explore a wide range of components, from databases and clusters to serverless functions.',
      },
      {
        number: 3,
        x: 88,
        y: 36,
        title: 'One-Click Deploy',
        description:
          'Deploy any component directly to your connected cloud account with a single click.',
      },
    ],
  },
  {
    step_number: 3,
    title: 'Cloud & Service Connections',
    description:
      'Securely link your cloud providers, source code repositories, and other services to the platform.',
    image_url:
      'https://qtrypzzcjebvfcihiynt.supabase.co/storage/v1/object/public/base44-prod/public/cdbc42ebd_Xnapper-2025-08-29-071716.png',
    callouts: [
      {
        number: 1,
        x: 50,
        y: 55,
        title: 'Provider Grid',
        description:
          'Connect to major clouds like AWS, Azure, and GCP, plus tools like GitHub and Docker.',
      },
      {
        number: 2,
        x: 23,
        y: 50,
        title: 'Guided Setup',
        description:
          "Click 'Connect' to start a simple, guided process for adding new account credentials securely.",
      },
    ],
  },
  {
    step_number: 4,
    title: 'IaC Module Registry',
    description:
      'Manage and version your reusable Infrastructure-as-Code modules for consistent and repeatable deployments.',
    image_url:
      'https://qtrypzzcjebvfcihiynt.supabase.co/storage/v1/object/public/base44-prod/public/f5a918138_Xnapper-2025-08-29-071725.png',
    callouts: [
      {
        number: 1,
        x: 40,
        y: 50,
        title: 'Versioned Modules',
        description:
          'View all your Terraform, Pulumi, and other IaC modules, complete with version history.',
      },
      {
        number: 2,
        x: 70,
        y: 45,
        title: 'Promote & Inspect',
        description:
          'Promote modules to different environments or jump directly to the source code for review.',
      },
    ],
  },
  {
    step_number: 5,
    title: 'Infra Hub - Visual Designer',
    description:
      'Visually map out and manage your cloud resources and their dependencies with an interactive graph.',
    image_url:
      'https://qtrypzzcjebvfcihiynt.supabase.co/storage/v1/object/public/base44-prod/public/1a6bb2f68_Xnapper-2025-08-29-071733.png',
    callouts: [
      {
        number: 1,
        x: 50,
        y: 60,
        title: 'Resource Graph',
        description:
          'See your entire infrastructure landscape, from VPCs to individual services, in a clear dependency graph.',
      },
      {
        number: 2,
        x: 40,
        y: 28,
        title: 'View Toggles',
        description:
          'Switch between a high-level view of your cloud inventory and specific infrastructure projects.',
      },
    ],
  },
  {
    step_number: 6,
    title: 'Service Hub Overview',
    description:
      'Organize and monitor all your microservices and their associated CI/CD pipelines from a central hub.',
    image_url:
      'https://qtrypzzcjebvfcihiynt.supabase.co/storage/v1/object/public/base44-prod/public/d51915c5e_Xnapper-2025-08-29-071751.png',
    callouts: [
      {
        number: 1,
        x: 35,
        y: 60,
        title: 'Recent Pipelines',
        description: 'Keep track of recent build and deployment pipelines for all your services.',
      },
      {
        number: 2,
        x: 75,
        y: 60,
        title: 'Services List',
        description:
          'View a list of all deployed services, with links to their source code repositories.',
      },
      {
        number: 3,
        x: 90,
        y: 40,
        title: 'Add New Service',
        description:
          'Onboard a new service and configure its deployment pipeline in just a few steps.',
      },
    ],
  },
  {
    step_number: 7,
    title: 'Source Code Template Library',
    description:
      'Accelerate development by using standardized templates for new services and CI/CD workflows.',
    image_url:
      'https://qtrypzzcjebvfcihiynt.supabase.co/storage/v1/object/public/base44-prod/public/21d5bbedc_Xnapper-2025-08-29-071801.png',
    callouts: [
      {
        number: 1,
        x: 40,
        y: 60,
        title: 'Template Library',
        description:
          'Choose from a library of templates for different languages and frameworks, like Next.js or Java.',
      },
      {
        number: 2,
        x: 80,
        y: 50,
        title: 'Filter by Type',
        description:
          'Filter templates by type, such as GitHub Actions, or maintainer to quickly find what you need.',
      },
    ],
  },
  {
    step_number: 8,
    title: 'Service Configuration',
    description:
      "Fine-tune every aspect of a service's deployment, from pipeline settings to build methods.",
    image_url:
      'https://qtrypzzcjebvfcihiynt.supabase.co/storage/v1/object/public/base44-prod/public/972cf266e_Xnapper-2025-08-29-071810.png',
    callouts: [
      {
        number: 1,
        x: 50,
        y: 50,
        title: 'Pipeline Configuration',
        description:
          'Configure how your service is built and deployed, choosing managed or self-managed pipelines.',
      },
      {
        number: 2,
        x: 50,
        y: 80,
        title: 'Deployment Rules',
        description:
          'Set up deployment triggers, pull request rules, and preview environments to match your workflow.',
      },
    ],
  },
  {
    step_number: 9,
    title: 'Visual Pipeline Builder',
    description: 'Visualize and manage the stages of your deployment pipelines in real-time.',
    image_url:
      'https://qtrypzzcjebvfcihiynt.supabase.co/storage/v1/object/public/base44-prod/public/9aa8d8023_Xnapper-2025-08-29-071818.png',
    callouts: [
      {
        number: 1,
        x: 20,
        y: 70,
        title: 'Pipeline Stages',
        description:
          'Track the progress of your pipeline through each stage, from build to deploy.',
      },
      {
        number: 2,
        x: 60,
        y: 70,
        title: 'Stage Details',
        description:
          'Click on a stage to see detailed steps and their execution status and duration.',
      },
      {
        number: 3,
        x: 60,
        y: 35,
        title: 'Run Summary',
        description:
          'Get a quick summary of the pipeline run, including trigger details, status, and total duration.',
      },
    ],
  },
];

const preloadImages = () => {
  TOUR_STEPS.forEach((step) => {
    const img = document.createElement('img');
    img.src = step.image_url;
  });
};

export default function TourPage() {
  const router = useRouter();
  const [currentStep, setCurrentStep] = useState(0);
  const [isStarted, setIsStarted] = useState(false);
  const [calloutsVisible, setCalloutsVisible] = useState(false);
  const [showDemoPrompt, setShowDemoPrompt] = useState(false);
  const [showClosingScreen, setShowClosingScreen] = useState(false);
  const [currentCalloutIndex, setCurrentCalloutIndex] = useState(0);

  useEffect(() => {
    preloadImages();
  }, []);

  // Reset callout index when step changes.
  useEffect(() => {
    const raf = requestAnimationFrame(() => setCurrentCalloutIndex(0));
    return () => cancelAnimationFrame(raf);
  }, [currentStep]);

  const shouldShowCallouts = isStarted && !showDemoPrompt && !showClosingScreen;

  useEffect(() => {
    if (shouldShowCallouts) {
      const timer = setTimeout(() => setCalloutsVisible(true), 500);
      return () => clearTimeout(timer);
    }
    const raf = requestAnimationFrame(() => setCalloutsVisible(false));
    return () => cancelAnimationFrame(raf);
  }, [currentStep, shouldShowCallouts]);

  const handleStart = () => {
    setIsStarted(true);
    setCalloutsVisible(false);
    setShowDemoPrompt(false);
    setShowClosingScreen(false);
  };

  const handleNext = () => {
    const currentStepData = TOUR_STEPS[currentStep];

    if (currentStepData.callouts && currentStepData.callouts.length > 0 && !calloutsVisible) {
      setCalloutsVisible(true);
      return;
    }

    if (currentStep === 3 && !showDemoPrompt) {
      setCalloutsVisible(false);
      setShowDemoPrompt(true);
      return;
    }
    if (currentStep === TOUR_STEPS.length - 1) {
      setCalloutsVisible(false);
      setIsStarted(false);
      setShowClosingScreen(true);
      return;
    }
    if (currentStep < TOUR_STEPS.length - 1) {
      setCalloutsVisible(false);
      setTimeout(() => setCurrentStep((prev) => prev + 1), 200);
    }
  };

  const handleKeepGoing = () => {
    setShowDemoPrompt(false);
    setTimeout(() => {
      setCurrentStep((prev) => prev + 1);
      setCalloutsVisible(true);
    }, 300);
  };

  const handleScheduleDemo = () => {
    setShowDemoPrompt(false);
    setShowClosingScreen(false);
    setIsStarted(false);
    setCurrentStep(0);
    router.push('/book-demo');
  };

  const handlePrevious = () => {
    if (currentStep > 0) {
      setCalloutsVisible(false);
      setTimeout(() => setCurrentStep((prev) => prev - 1), 200);
    }
  };

  const handleClosingContinue = () => {
    setShowClosingScreen(false);
    setIsStarted(false);
    setCurrentStep(0);
  };

  const handleCalloutNext = () => {
    const currentStepData = TOUR_STEPS[currentStep];
    if (currentStepData.callouts && currentCalloutIndex < currentStepData.callouts.length - 1) {
      setCurrentCalloutIndex((prev) => prev + 1);
    } else {
      handleNext();
    }
  };

  const handleCalloutPrevious = () => {
    if (currentCalloutIndex > 0) {
      setCurrentCalloutIndex((prev) => prev - 1);
    }
  };

  const progress = ((currentStep + 1) / TOUR_STEPS.length) * 100;

  if (showClosingScreen) {
    return (
      <div className="min-h-screen bg-[#0a0a0a]">
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          className="fixed inset-0 z-[9998] bg-black/60 backdrop-blur-sm flex items-center justify-center p-4"
          style={{ isolation: 'isolate' }}
        >
          <div className="bg-white rounded-2xl shadow-2xl p-8 max-w-md w-full text-center">
            <div className="w-16 h-16 bg-[#111] rounded-full flex items-center justify-center mx-auto mb-6">
              <svg
                className="w-8 h-8 text-white"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M5 13l4 4L19 7"
                />
              </svg>
            </div>
            <h2 className="text-2xl font-bold text-gray-900 mb-4">
              Thank you for taking the interactive tour!
            </h2>
            <p className="text-gray-600 mb-8 leading-relaxed">
              If you want to learn more and interact with one of our engineers, please book a demo
              to see how our platform can transform your DevOps workflows.
            </p>
            <div className="space-y-3">
              <Button
                onClick={handleScheduleDemo}
                className="w-full bg-[#111] hover:bg-[#222] text-white py-3 rounded-xl"
              >
                Book Demo
              </Button>
              <Button
                variant="ghost"
                onClick={handleClosingContinue}
                className="w-full text-gray-500 hover:text-gray-700"
              >
                Continue to Site
              </Button>
            </div>
          </div>
        </motion.div>
      </div>
    );
  }

  if (!isStarted) {
    return (
      <div className="min-h-screen bg-[#0a0a0a]">
        <Toaster />
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          className="fixed inset-0 z-[9997] bg-black/60 backdrop-blur-sm flex items-center justify-center p-4"
          style={{ isolation: 'isolate' }}
        >
          <div className="bg-white rounded-2xl shadow-2xl p-8 max-w-md w-full text-center">
            <div className="w-16 h-16 bg-[#111] rounded-full flex items-center justify-center mx-auto mb-6">
              <Play className="w-8 h-8 text-white ml-1" />
            </div>
            <h2 className="text-2xl font-bold text-gray-900 mb-4">Welcome to the Platform Tour</h2>
            <p className="text-gray-600 mb-8 leading-relaxed">
              Discover how our DevOps platform simplifies infrastructure management, deployment, and
              monitoring with powerful visual tools and integrations.
            </p>
            <Button
              onClick={handleStart}
              className="w-full bg-[#111] hover:bg-[#222] text-white py-3 rounded-xl"
            >
              Start Interactive Tour
            </Button>
          </div>
        </motion.div>
      </div>
    );
  }

  const currentStepData = TOUR_STEPS[currentStep];

  return (
    <div
      className="min-h-screen bg-[#0a0a0a]"
      style={{ isolation: 'isolate' }}
    >
      <Toaster />
      <AnimatePresence>
        {showDemoPrompt && (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            className="fixed inset-0 z-[9999] bg-black/60 backdrop-blur-sm flex items-center justify-center p-4"
            style={{ isolation: 'isolate' }}
          >
            <motion.div
              initial={{ opacity: 0, scale: 0.9, y: 20 }}
              animate={{ opacity: 1, scale: 1, y: 0 }}
              exit={{ opacity: 0, scale: 0.9, y: 20 }}
              className="bg-white rounded-2xl shadow-2xl p-8 max-w-lg w-full relative"
            >
              <button
                onClick={() => setShowDemoPrompt(false)}
                className="absolute top-4 right-4 text-gray-400 hover:text-gray-600 transition-colors"
              >
                <X className="w-5 h-5" />
              </button>
              <div className="text-center mb-6">
                <div className="w-12 h-12 bg-[#111] rounded-full flex items-center justify-center mx-auto mb-4">
                  <Play className="w-6 h-6 text-white ml-1" />
                </div>
                <h3 className="text-2xl font-bold text-gray-900 mb-4">Want to see more?</h3>
                <p className="text-gray-600 mb-4 leading-relaxed">
                  Keep going to see how our platform can give engineering teams more visibility into
                  infrastructure deployment and monitoring workflows.
                </p>
                <p className="text-gray-600 leading-relaxed">
                  Or, schedule a demo with our team to see it live and get all of your questions
                  answered.
                </p>
              </div>
              <div className="flex flex-col sm:flex-row gap-3">
                <Button
                  variant="outline"
                  onClick={handleKeepGoing}
                  className="flex-1 border-2 border-gray-200 hover:border-gray-300 hover:bg-gray-50 text-gray-700"
                >
                  Keep going
                </Button>
                <Button
                  onClick={handleScheduleDemo}
                  className="flex-1 bg-[#111] hover:bg-[#222] text-white"
                >
                  Schedule demo
                </Button>
              </div>
            </motion.div>
          </motion.div>
        )}
      </AnimatePresence>

      <div className="relative w-full h-screen p-4 md:p-8 flex items-center justify-center">
        <div className="max-w-6xl w-full h-[calc(100vh-160px)]">
          <AnimatePresence mode="wait">
            <motion.div
              key={currentStep}
              initial={{ opacity: 0, scale: 0.95 }}
              animate={{ opacity: 1, scale: 1 }}
              exit={{ opacity: 0, scale: 0.95 }}
              transition={{ duration: 0.4 }}
              className="relative w-full h-full bg-gray-100 rounded-xl overflow-hidden shadow-2xl"
            >
              <Image
                src={currentStepData.image_url}
                alt={currentStepData.title}
                fill
                sizes="100vw"
                className="object-contain"
              />
              <div className="absolute inset-0 bg-black/10" />
              {currentStepData.callouts?.map((callout, index) => (
                <CalloutCard
                  key={index}
                  callout={callout}
                  index={index}
                  isVisible={calloutsVisible}
                  isCurrent={index === currentCalloutIndex}
                  currentCalloutIndex={currentCalloutIndex}
                  totalCallouts={currentStepData.callouts?.length || 0}
                  onNumberClick={(index) => {
                    setCurrentCalloutIndex(index);
                    if (!calloutsVisible) {
                      setCalloutsVisible(true);
                    }
                  }}
                  onNext={handleCalloutNext}
                  onPrevious={handleCalloutPrevious}
                />
              ))}
            </motion.div>

            <motion.div
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: 0.3 }}
              className="absolute top-20 left-6 bg-white/95 backdrop-blur-sm rounded-xl p-4 shadow-lg max-w-xs z-10"
            >
              <h3 className="font-bold text-gray-900 mb-2">{currentStepData.title}</h3>
              <p className="text-sm text-gray-600">{currentStepData.description}</p>
            </motion.div>
          </AnimatePresence>
        </div>
      </div>

      {isStarted && (
        <motion.div
          initial={{ opacity: 0, y: 50 }}
          animate={{ opacity: 1, y: 0 }}
          className="fixed bottom-0 left-0 right-0 z-[9996] bg-white/95 backdrop-blur-sm border-t border-gray-200 p-6"
        >
          <div className="max-w-4xl mx-auto">
            <div className="flex justify-between items-center">
              <Button
                variant="outline"
                onClick={handlePrevious}
                disabled={currentStep === 0}
                className="flex items-center gap-2 h-10"
              >
                <ChevronLeft className="w-4 h-4" />
                Previous
              </Button>
              <div className="flex-1 mx-8">
                <div className="flex justify-between items-center mb-2">
                  <span className="text-sm font-medium text-gray-700">
                    Step {currentStep + 1} of {TOUR_STEPS.length}
                  </span>
                  <span className="text-sm text-gray-500">{Math.round(progress)}% Complete</span>
                </div>
                <div className="w-full bg-gray-200 rounded-full h-2">
                  <motion.div
                    initial={{ width: 0 }}
                    animate={{ width: `${progress}%` }}
                    transition={{ duration: 0.5, ease: 'easeOut' }}
                    className="bg-white h-2 rounded-full"
                  />
                </div>
              </div>
              <Button
                onClick={handleNext}
                className="flex items-center gap-2 bg-[#111] hover:bg-[#222] h-10"
              >
                {currentStep === TOUR_STEPS.length - 1 ? 'Finish Tour' : 'Next'}
                <ChevronRight className="w-4 h-4" />
              </Button>
            </div>
          </div>
        </motion.div>
      )}
    </div>
  );
}
