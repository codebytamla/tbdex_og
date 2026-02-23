import React, { useState } from 'react';
import { Moon, Sun, RefreshCw, ChevronRight, Home, Users, Activity, Plus } from 'lucide-react';

export default function App() {
  const [isDarkMode, setIsDarkMode] = useState(true);

  const tasks = [
    { id: 1, name: 'John Mwangi', location: 'Kipkelion', time: '9:00 AM', status: 'Pending' },
    { id: 2, name: 'Esther Kimani', location: 'Eldoret', time: '10:30 AM', status: 'In Progress' },
    { id: 3, name: 'David Omondi', location: 'Nakuru', time: '1:00 PM', status: 'Done' },
    { id: 4, name: 'Grace Achieng', location: 'Kericho', time: '2:30 PM', status: 'Pending' },
    { id: 5, name: 'James Otieno', location: 'Kipkelion', time: '3:45 PM', status: 'Pending' },
  ];

  const getStatusBadge = (status) => {
    switch (status) {
      case 'Pending':
        return 'bg-amber-100 dark:bg-amber-900/40 text-amber-800 dark:text-amber-200 border border-amber-200 dark:border-amber-800/60';
      case 'In Progress':
        return 'bg-blue-100 dark:bg-blue-900/40 text-blue-800 dark:text-blue-200 border border-blue-200 dark:border-blue-800/60';
      case 'Done':
        return 'bg-emerald-100 dark:bg-emerald-900/40 text-emerald-800 dark:text-emerald-200 border border-emerald-200 dark:border-emerald-800/60';
      default:
        return 'bg-gray-100 dark:bg-neutral-800 text-gray-800 dark:text-neutral-200 border border-gray-200 dark:border-neutral-700';
    }
  };

  return (
    <div className={isDarkMode ? 'dark' : ''}>
      <div className="flex justify-center bg-gray-100 dark:bg-neutral-950 min-h-screen font-sans transition-colors duration-300">
        {/* Mobile Device Container */}
        <div className="w-full max-w-md bg-[#FAFAFC] dark:bg-[#131314] min-h-screen shadow-2xl relative pb-24 overflow-hidden transition-colors duration-300">
          
          {/* Header */}
          <header className="px-6 pt-12 pb-4 flex justify-between items-start">
            <div>
              <p className="text-slate-500 dark:text-neutral-400 text-sm mb-1 transition-colors duration-300">Today</p>
              <h1 className="text-2xl font-bold text-slate-900 dark:text-neutral-100 tracking-tight transition-colors duration-300">Task Overview</h1>
            </div>
            <div className="flex items-center gap-4">
              <button 
                onClick={() => setIsDarkMode(!isDarkMode)}
                className="text-slate-600 dark:text-neutral-400 hover:text-slate-900 dark:hover:text-neutral-100 transition-colors duration-300"
              >
                {isDarkMode ? <Sun size={22} className="fill-current" /> : <Moon size={22} className="fill-current" />}
              </button>
              <div className="w-10 h-10 rounded-full bg-blue-100 dark:bg-blue-900/50 text-blue-700 dark:text-blue-200 flex items-center justify-center font-semibold text-sm border border-blue-200 dark:border-blue-800/60 transition-colors duration-300">
                TD
              </div>
            </div>
          </header>

          {/* Content Area */}
          <main className="px-6 flex flex-col gap-6">
            
            {/* Sync Banner */}
            <button className="w-full bg-amber-50 dark:bg-amber-900/20 border border-amber-200 dark:border-amber-800/50 rounded-xl p-3 flex items-center justify-between text-amber-800 dark:text-amber-300 shadow-sm active:scale-[0.98] transition-all duration-300">
              <div className="flex items-center gap-3 text-sm font-medium">
                <RefreshCw size={16} />
                <span>Last synced: 2 hours ago</span>
              </div>
              <ChevronRight size={18} className="text-amber-600 dark:text-amber-500" />
            </button>

            {/* Metrics Grid */}
            <div className="grid grid-cols-3 gap-3">
              {/* Pending */}
              <div className="bg-blue-50 dark:bg-blue-900/20 border border-blue-100 dark:border-blue-800/50 rounded-2xl p-4 flex flex-col items-center justify-center text-center shadow-sm transition-colors duration-300">
                <span className="text-2xl font-bold text-blue-700 dark:text-blue-400 mb-1">3</span>
                <span className="text-xs font-medium text-blue-800 dark:text-blue-300 leading-tight">Pending<br/>Assessments</span>
              </div>
              {/* Referrals */}
              <div className="bg-rose-50 dark:bg-rose-900/20 border border-rose-100 dark:border-rose-800/50 rounded-2xl p-4 flex flex-col items-center justify-center text-center shadow-sm transition-colors duration-300">
                <span className="text-2xl font-bold text-rose-700 dark:text-rose-400 mb-1">3</span>
                <span className="text-xs font-medium text-rose-800 dark:text-rose-300 leading-tight">Referrals<br/>Today</span>
              </div>
              {/* Total */}
              <div className="bg-emerald-50 dark:bg-emerald-900/20 border border-emerald-100 dark:border-emerald-800/50 rounded-2xl p-4 flex flex-col items-center justify-center text-center shadow-sm transition-colors duration-300">
                <span className="text-2xl font-bold text-emerald-700 dark:text-emerald-400 mb-1">6</span>
                <span className="text-xs font-medium text-emerald-800 dark:text-emerald-300 leading-tight">Total<br/>Patients</span>
              </div>
            </div>

            {/* Recent Assessments */}
        <div className="mt-4">
          <h2 className="text-lg font-bold text-slate-900 dark:text-neutral-100 mb-4 transition-colors duration-300">Recent Assessments</h2>
          <div className="flex flex-col gap-3">
            {tasks.map((task) => (
              <div key={task.id} className="bg-white dark:bg-[#1E1F20] border border-gray-100 dark:border-neutral-800 rounded-2xl p-4 flex items-center justify-between shadow-sm transition-colors duration-300">
                <div>
                  <h3 className="font-semibold text-slate-900 dark:text-neutral-100 mb-1">{task.name}</h3>
                  <p className="text-xs text-slate-500 dark:text-neutral-400 font-medium">
                    {task.location} • {task.time}
                  </p>
                </div>
                <span className={`px-3 py-1 rounded-full text-xs font-semibold tracking-wide transition-colors duration-300 ${getStatusBadge(task.status)}`}>
                  {task.status}
                </span>
              </div>
            ))}
          </div>
        </div>
      </main>

      {/* Bottom Navigation & FAB */}
      <div className="absolute bottom-6 left-6 right-6 flex items-center gap-4 z-10">
        
        {/* Nav Capsule */}
        <div className="flex-1 bg-white/90 dark:bg-[#1E1F20]/90 backdrop-blur-md border border-gray-100 dark:border-neutral-800 rounded-full p-2 flex items-center justify-between shadow-xl shadow-slate-200/50 dark:shadow-black/50 transition-colors duration-300">
          <button className="flex flex-col items-center justify-center w-1/3 py-2 bg-blue-100 dark:bg-blue-900/30 text-blue-700 dark:text-blue-400 rounded-full transition-colors duration-300">
            <Home size={20} className="mb-1" />
            <span className="text-[10px] font-semibold">Home</span>
          </button>
          <button className="flex flex-col items-center justify-center w-1/3 py-2 text-slate-400 dark:text-neutral-500 hover:text-slate-600 dark:hover:text-neutral-200 transition-colors duration-300">
            <Users size={20} className="mb-1" />
            <span className="text-[10px] font-semibold">Patients</span>
          </button>
          <button className="flex flex-col items-center justify-center w-1/3 py-2 text-slate-400 dark:text-neutral-500 hover:text-slate-600 dark:hover:text-neutral-200 transition-colors duration-300">
            <Activity size={20} className="mb-1" />
            <span className="text-[10px] font-semibold">Activity</span>
          </button>
        </div>
        
        {/* Floating Action Button */}
            <button className="w-16 h-16 shrink-0 bg-blue-600 dark:bg-blue-500 text-white rounded-full flex items-center justify-center shadow-xl shadow-blue-600/30 dark:shadow-blue-900/40 hover:bg-blue-700 dark:hover:bg-blue-600 active:scale-95 transition-all">
              <Plus size={28} />
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
