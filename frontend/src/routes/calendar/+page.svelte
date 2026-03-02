<script lang="ts">
import { onMount } from 'svelte';
import { goto } from '$app/navigation';
import { taskAPI } from '$lib/api';
import { wsService } from '$lib/websocket';
import {
		Calendar as CalendarIcon, ChevronLeft, ChevronRight, Circle,
		Plus, Edit2, Trash2, X, Bell, Star,
		Archive, Download, RotateCcw, CreditCard, AlertCircle, Check,
		History
	} from 'lucide-svelte';

	// ── Types ──────────────────────────────────────────────────────────────────
type TaskStatus    = 'todo' | 'in-progress' | 'done';
type TaskPriority  = 'low' | 'medium' | 'high';
type CalendarSub   = 'regular' | 'payment' | 'reminder' | 'event';
type RecurringType = 'daily' | 'weekly' | 'monthly' | 'yearly';
type AppTab        = 'calendar' | 'history' | 'archive';
type HistoryFilter = 'all' | 'paid' | 'unpaid' | 'overdue';

type Task = {
id: number;
title: string;
description?: string | null;
priority: TaskPriority;
status: TaskStatus;
due_date?: string | null;
tags?: string | null;
is_quick_task?: boolean;
is_archived?: boolean;
quadrant?: string | null;
task_type?: string;
is_payment?: boolean;
amount?: number;
currency?: string;
is_paid?: boolean;
paid_at?: string | null;
payment_notes?: string | null;
calendar_subtype?: CalendarSub;
created_at?: string;
updated_at?: string;
};

type NewTask = {
title: string;
description: string;
priority: TaskPriority;
status: TaskStatus;
tags: string;
due_date?: string;
is_quick_task: boolean;
quadrant: string;
task_type: string;
calendar_subtype: CalendarSub;
is_payment: boolean;
amount: number;
currency: string;
is_paid: boolean;
payment_notes: string;
};

// ── State ──────────────────────────────────────────────────────────────────
let tasks         = $state<Task[]>([]);
let archivedTasks = $state<Task[]>([]);
let loading       = $state(true);

let currentMonth = $state(new Date().getMonth());
let currentYear  = $state(new Date().getFullYear());
let selectedDate = $state<number | null>(null);

let showAddTask   = $state(false);
let editingTask   = $state<Task | null>(null);
let showMarkPaid  = $state<Task | null>(null);

let activeTab       = $state<AppTab>('calendar');
let historyFilter   = $state<HistoryFilter>('all');
let archiveFilter   = $state<'all' | 'payment' | 'event' | 'reminder'>('all');
let archiveSearch   = $state('');

let recurring      = $state(false);
let recurringType  = $state<RecurringType>('weekly');
let recurringCount = $state(4);

let newTask = $state<NewTask>({
title: '', description: '', priority: 'medium', status: 'todo',
tags: '', due_date: '', is_quick_task: false, quadrant: '',
task_type: 'calendar', calendar_subtype: 'regular',
is_payment: false, amount: 0, currency: 'USD', is_paid: false, payment_notes: ''
});

// ── Constants ──────────────────────────────────────────────────────────────
const monthNames    = ['January','February','March','April','May','June',
                       'July','August','September','October','November','December'];
const dayNames      = ['Sun','Mon','Tue','Wed','Thu','Fri','Sat'];
const dayNamesShort = ['S','M','T','W','T','F','S'];
const currencies    = ['USD','EUR','GBP','INR','AUD','CAD','SGD','AED','MYR','PHP','THB','JPY'];

// ── Lifecycle ──────────────────────────────────────────────────────────────
onMount(async () => {
const token = typeof window !== 'undefined' ? localStorage.getItem('authToken') : null;
if (!token) { goto('/login'); return; }
try {
await Promise.all([loadTasks(), loadArchivedTasks()]);
} catch (e: any) {
console.error(e);
if (e.message?.includes('401') || e.message?.includes('Authorization')) {
localStorage.removeItem('authToken'); goto('/login');
}
} finally { loading = false; }

const refresh = async () => { await Promise.all([loadTasks(), loadArchivedTasks()]); };
wsService.on('task_create', refresh);
wsService.on('task_update', refresh);
wsService.on('task_delete', refresh);
return () => {
wsService.off('task_create', refresh);
wsService.off('task_update', refresh);
wsService.off('task_delete', refresh);
};
});

// ── Loaders ───────────────────────────────────────────────────────────────
async function loadTasks() {
const data = await taskAPI.getAll();
tasks = data.filter((t: Task) => !t.is_archived && t.due_date && t.task_type === 'calendar');
}

async function loadArchivedTasks() {
try {
const data = await taskAPI.getArchived();
archivedTasks = data.filter((t: Task) => t.task_type === 'calendar');
} catch { archivedTasks = []; }
}

// ── Calendar helpers ──────────────────────────────────────────────────────
function getDaysInMonth(m: number, y: number)  { return new Date(y, m + 1, 0).getDate(); }
function getFirstDay(m: number, y: number)     { return new Date(y, m, 1).getDay(); }

function isDateToday(day: number) {
const t = new Date();
return day === t.getDate() && currentMonth === t.getMonth() && currentYear === t.getFullYear();
}
function isDatePast(day: number) {
const d = new Date(currentYear, currentMonth, day);
const t = new Date(); t.setHours(0,0,0,0);
return d < t;
}
function getTasksForDate(day: number): Task[] {
const ds = `${currentYear}-${String(currentMonth+1).padStart(2,'0')}-${String(day).padStart(2,'0')}`;
return tasks.filter(t => t.due_date?.startsWith(ds));
}
function getCalendarDays(): (number | null)[] {
const days: (number | null)[] = [];
for (let i = 0; i < getFirstDay(currentMonth, currentYear); i++) days.push(null);
for (let d = 1; d <= getDaysInMonth(currentMonth, currentYear); d++) days.push(d);
return days;
}
function previousMonth() {
if (currentMonth === 0) { currentMonth = 11; currentYear--; } else { currentMonth--; }
selectedDate = null;
}
function nextMonth() {
if (currentMonth === 11) { currentMonth = 0; currentYear++; } else { currentMonth++; }
selectedDate = null;
}
function goToToday() {
const t = new Date(); currentMonth = t.getMonth(); currentYear = t.getFullYear();
selectedDate = t.getDate();
}

// ── Styling helpers ───────────────────────────────────────────────────────
function getPriorityBorder(p: TaskPriority) {
if (p === 'high')   return 'border-l-red-500 bg-red-950/30';
if (p === 'medium') return 'border-l-yellow-500 bg-yellow-950/20';
return 'border-l-gray-600 bg-gray-800/50';
}
function getDotColor(task: Task, isToday: boolean): string {
if (isToday) return 'bg-white';
const sub = task.calendar_subtype || 'regular';
if (sub === 'payment' || task.is_payment) {
if (task.is_paid) return 'bg-green-400';
const past = task.due_date && new Date(task.due_date) < new Date();
return past ? 'bg-red-500' : 'bg-amber-400';
}
if (sub === 'reminder') return 'bg-purple-400';
if (sub === 'event')    return 'bg-blue-400';
if (task.status === 'done') return 'bg-green-500';
if (task.status === 'in-progress') return 'bg-blue-400';
return 'bg-gray-500';
}
function getSubtypeBadge(sub?: CalendarSub) {
if (sub === 'payment')  return 'bg-amber-950 text-amber-300 border-amber-800';
if (sub === 'reminder') return 'bg-purple-950 text-purple-300 border-purple-800';
if (sub === 'event')    return 'bg-blue-950 text-blue-300 border-blue-800';
return 'bg-gray-800 text-gray-400 border-gray-700';
}
function subIcon(sub?: CalendarSub) {
if (sub === 'payment')  return CreditCard;
if (sub === 'reminder') return Bell;
if (sub === 'event')    return Star;
return Circle;
}
function subIconColor(task: Task): string {
const sub = task.calendar_subtype || 'regular';
if (sub === 'payment' || task.is_payment) {
if (task.is_paid) return 'text-green-400';
const past = task.due_date && new Date(task.due_date) < new Date();
return past ? 'text-red-400' : 'text-amber-400';
}
if (sub === 'reminder') return 'text-purple-400';
if (sub === 'event')    return 'text-blue-400';
if (task.status === 'done') return 'text-green-400';
return 'text-gray-400';
}

function isOverdue(task: Task): boolean {
if (!task.due_date || task.is_paid || task.status === 'done') return false;
return new Date(task.due_date) < new Date();
}

function formatCurrency(amount: number, currency: string): string {
try { return new Intl.NumberFormat('en-US', { style: 'currency', currency }).format(amount); }
catch { return `${currency} ${amount.toFixed(2)}`; }
}
function formatDate(d?: string | null): string {
if (!d) return '—';
return new Date(d).toLocaleDateString('en-US', { year: 'numeric', month: 'short', day: 'numeric' });
}
function groupByMonth(list: Task[]): Map<string, Task[]> {
const map = new Map<string, Task[]>();
for (const t of list) {
const key = t.due_date
? new Date(t.due_date).toLocaleDateString('en-US', { year:'numeric', month:'long' })
: 'No Date';
if (!map.has(key)) map.set(key, []);
map.get(key)!.push(t);
}
return map;
}

// ── Derived ───────────────────────────────────────────────────────────────
const calendarDays       = $derived(getCalendarDays());
const selectedDateTasks  = $derived(selectedDate ? getTasksForDate(selectedDate) : []);

const paymentStats = $derived((() => {
const mp = tasks.filter(t => {
if (!t.is_payment && t.calendar_subtype !== 'payment') return false;
if (!t.due_date) return false;
const d = new Date(t.due_date);
return d.getMonth() === currentMonth && d.getFullYear() === currentYear;
});
const totalAmt = mp.reduce((s, t) => s + (t.amount || 0), 0);
const paid     = mp.filter(t => t.is_paid);
const paidAmt  = paid.reduce((s, t) => s + (t.amount || 0), 0);
const overdue  = mp.filter(t => isOverdue(t));
return { count: mp.length, totalAmt, paidCount: paid.length, paidAmt,
         unpaidAmt: totalAmt - paidAmt, overdueCount: overdue.length,
         currency: mp[0]?.currency || 'USD' };
})());

const allPaymentTasks = $derived((() => {
const all = [...tasks, ...archivedTasks].filter(t => t.is_payment || t.calendar_subtype === 'payment');
all.sort((a, b) => new Date(b.due_date || 0).getTime() - new Date(a.due_date || 0).getTime());
if (historyFilter === 'paid')    return all.filter(t => t.is_paid);
if (historyFilter === 'unpaid')  return all.filter(t => !t.is_paid);
if (historyFilter === 'overdue') return all.filter(t => isOverdue(t));
return all;
})());

const filteredArchive = $derived((() => {
let list = [...archivedTasks];
if (archiveFilter !== 'all') {
if (archiveFilter === 'payment') list = list.filter(t => t.calendar_subtype === 'payment' || t.is_payment);
else list = list.filter(t => (t.calendar_subtype || 'regular') === archiveFilter);
}
if (archiveSearch.trim()) {
const q = archiveSearch.toLowerCase();
list = list.filter(t => t.title.toLowerCase().includes(q) || (t.description || '').toLowerCase().includes(q));
}
return list.sort((a, b) => new Date(b.updated_at || 0).getTime() - new Date(a.updated_at || 0).getTime());
})());

// ── Form ──────────────────────────────────────────────────────────────────
function openAddTaskModal(day?: number) {
const d = day || selectedDate;
if (d) newTask.due_date = `${currentYear}-${String(currentMonth+1).padStart(2,'0')}-${String(d).padStart(2,'0')}`;
showAddTask = true;
}
function closeAddTaskModal() { showAddTask = false; resetNewTask(); }
function resetNewTask() {
newTask = { title:'', description:'', priority:'medium', status:'todo', tags:'', due_date:'',
is_quick_task:false, quadrant:'', task_type:'calendar', calendar_subtype:'regular',
is_payment:false, amount:0, currency:'USD', is_paid:false, payment_notes:'' };
recurring = false; recurringType = 'weekly'; recurringCount = 4;
}
function computeRecurringDates(base: string, type: RecurringType, count: number): string[] {
const dates: string[] = [];
const b = new Date(base + 'T00:00:00Z');
for (let i = 0; i < count; i++) {
const d = new Date(b);
if (type === 'daily')        d.setUTCDate(d.getUTCDate() + i);
else if (type === 'weekly')  d.setUTCDate(d.getUTCDate() + i * 7);
else if (type === 'monthly') d.setUTCMonth(d.getUTCMonth() + i);
else if (type === 'yearly')  d.setUTCFullYear(d.getUTCFullYear() + i);
dates.push(d.toISOString());
}
return dates;
}

// ── CRUD ──────────────────────────────────────────────────────────────────
async function handleAddTask() {
if (!newTask.title.trim()) return;
try {
const isPayment = newTask.calendar_subtype === 'payment';
let dueDates: (string | undefined)[];
if (recurring && newTask.due_date) {
dueDates = computeRecurringDates(newTask.due_date, recurringType, recurringCount);
} else {
dueDates = [newTask.due_date ? new Date(newTask.due_date + 'T00:00:00Z').toISOString() : undefined];
}
for (const dd of dueDates) {
await taskAPI.create({
title: newTask.title, description: newTask.description,
priority: newTask.priority, status: newTask.status,
tags: newTask.tags, due_date: dd,
is_quick_task: false, quadrant: '',
task_type: 'calendar', calendar_subtype: newTask.calendar_subtype,
is_payment: isPayment, amount: isPayment ? newTask.amount : 0,
currency: newTask.currency, is_paid: newTask.is_paid,
payment_notes: newTask.payment_notes
});
}
await Promise.all([loadTasks(), loadArchivedTasks()]);
closeAddTaskModal();
} catch (e) { console.error(e); alert('Failed to create.'); }
}

async function handleSaveEdit() {
if (!editingTask || !editingTask.title.trim()) return;
try {
let fd: string | undefined;
if (editingTask.due_date) fd = new Date(editingTask.due_date.split('T')[0] + 'T00:00:00Z').toISOString();
const isPayment = editingTask.calendar_subtype === 'payment' || !!editingTask.is_payment;
await taskAPI.update(editingTask.id, {
title: editingTask.title, description: editingTask.description || '',
priority: editingTask.priority, status: editingTask.status,
tags: editingTask.tags || '', due_date: fd,
is_quick_task: false, quadrant: '', task_type: 'calendar',
calendar_subtype: editingTask.calendar_subtype || 'regular',
is_payment: isPayment, amount: editingTask.amount || 0,
currency: editingTask.currency || 'USD',
is_paid: !!editingTask.is_paid,
payment_notes: editingTask.payment_notes || ''
});
await Promise.all([loadTasks(), loadArchivedTasks()]);
editingTask = null;
} catch (e) { console.error(e); alert('Failed to save.'); }
}

async function handleDelete(id: number) {
if (!confirm('Delete this task?')) return;
try { await taskAPI.delete(id); await Promise.all([loadTasks(), loadArchivedTasks()]); }
catch (e) { console.error(e); alert('Failed to delete.'); }
}
async function handleArchive(id: number) {
try { await taskAPI.archive(id); await Promise.all([loadTasks(), loadArchivedTasks()]); }
catch { alert('Failed to archive.'); }
}
async function handleRestore(id: number) {
try { await taskAPI.restore(id); await Promise.all([loadTasks(), loadArchivedTasks()]); }
catch { alert('Failed to restore.'); }
}
async function handlePermDelete(id: number) {
if (!confirm('Permanently delete? Cannot be undone.')) return;
try { await taskAPI.permanentDelete(id); await loadArchivedTasks(); }
catch { alert('Failed to permanently delete.'); }
}

async function markAsPaid(task: Task) {
try {
const baseDate = task.due_date ? task.due_date.split('T')[0] : null;
await taskAPI.update(task.id, {
...task,
is_paid: true,
paid_at: new Date().toISOString(),
status: 'done',
due_date: baseDate ? new Date(baseDate + 'T00:00:00Z').toISOString() : task.due_date
});
await Promise.all([loadTasks(), loadArchivedTasks()]);
showMarkPaid = null;
} catch { alert('Failed to mark as paid.'); }
}
async function markAsUnpaid(task: Task) {
try {
const baseDate = task.due_date ? task.due_date.split('T')[0] : null;
await taskAPI.update(task.id, {
...task,
is_paid: false,
paid_at: null,
status: 'todo',
due_date: baseDate ? new Date(baseDate + 'T00:00:00Z').toISOString() : task.due_date
});
await Promise.all([loadTasks(), loadArchivedTasks()]);
} catch { alert('Failed to update.'); }
}

// ── CSV download ──────────────────────────────────────────────────────────
function downloadCSV(list: Task[], filename: string) {
const header = ['Date','Title','Description','Type','Amount','Currency','Priority','Status','Paid','Paid At','Notes','Created'];
const rows = list.map(t => [
t.due_date ? new Date(t.due_date).toLocaleDateString() : '',
`"${(t.title||'').replace(/"/g,'""')}"`,
`"${(t.description||'').replace(/"/g,'""')}"`,
t.calendar_subtype || 'regular',
t.amount || '',
t.currency || '',
t.priority || '',
t.status || '',
t.is_paid ? 'Yes' : 'No',
t.paid_at ? new Date(t.paid_at).toLocaleDateString() : '',
`"${(t.payment_notes||'').replace(/"/g,'""')}"`,
t.created_at ? new Date(t.created_at).toLocaleDateString() : ''
]);
const csv  = [header.join(','), ...rows.map(r => r.join(','))].join('\n');
const blob = new Blob([csv], { type: 'text/csv;charset=utf-8;' });
const url  = URL.createObjectURL(blob);
const a = document.createElement('a'); a.href = url; a.download = filename; a.click();
URL.revokeObjectURL(url);
}
</script>

<svelte:head><title>Calendar - Tonish</title></svelte:head>

<div class="max-w-[1600px] mx-auto">
{#if loading}
<div class="text-center py-12"><p class="text-gray-400 text-sm animate-pulse">Loading calendar...</p></div>
{:else}

<!-- ── Tab Bar ──────────────────────────────────────────────────────────── -->
<div class="flex items-center gap-1 mb-3 bg-gray-900 border border-gray-800 rounded-lg p-1 w-fit">
<button onclick={() => activeTab = 'calendar'}
class="flex items-center gap-1.5 px-3 py-1.5 rounded-md text-xs font-medium transition
{activeTab === 'calendar' ? 'bg-blue-600 text-white' : 'text-gray-400 hover:text-white hover:bg-gray-800'}">
<CalendarIcon size={13} /> Calendar
</button>
<button onclick={() => activeTab = 'history'}
class="flex items-center gap-1.5 px-3 py-1.5 rounded-md text-xs font-medium transition
{activeTab === 'history' ? 'bg-blue-600 text-white' : 'text-gray-400 hover:text-white hover:bg-gray-800'}">
<History size={13} /> Payment History
{#if allPaymentTasks.length > 0}
<span class="bg-amber-600 text-white text-[9px] rounded-full px-1.5 py-0.5 leading-none">{allPaymentTasks.length}</span>
{/if}
</button>
<button onclick={() => activeTab = 'archive'}
class="flex items-center gap-1.5 px-3 py-1.5 rounded-md text-xs font-medium transition
{activeTab === 'archive' ? 'bg-blue-600 text-white' : 'text-gray-400 hover:text-white hover:bg-gray-800'}">
<Archive size={13} /> Archive
{#if archivedTasks.length > 0}
<span class="bg-gray-600 text-white text-[9px] rounded-full px-1.5 py-0.5 leading-none">{archivedTasks.length}</span>
{/if}
</button>
</div>

<!-- ════════════════════ CALENDAR TAB ════════════════════════════════════ -->
{#if activeTab === 'calendar'}
<div class="grid grid-cols-1 sm:grid-cols-3 gap-2 sm:gap-3">

<!-- Calendar Grid Panel -->
<div class="sm:col-span-2 bg-gray-900 border border-gray-800 rounded-lg p-2.5 sm:p-3">
<div class="flex items-center justify-between mb-2.5">
<h2 class="text-sm font-bold text-white">
{monthNames[currentMonth]} <span class="text-gray-400 font-normal">{currentYear}</span>
</h2>
<div class="flex items-center gap-1">
<button onclick={() => openAddTaskModal()}
class="h-8 px-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 transition flex items-center gap-1 text-xs font-medium touch-manipulation">
<Plus size={12} /><span class="hidden sm:inline">Add</span>
</button>
<button onclick={goToToday}
class="h-8 px-2 text-xs bg-blue-950 text-blue-300 rounded-md hover:bg-blue-900 font-medium transition touch-manipulation">Today</button>
<button onclick={previousMonth}
class="w-8 h-8 flex items-center justify-center hover:bg-gray-800 rounded-md text-gray-300 touch-manipulation">
<ChevronLeft size={15} /></button>
<button onclick={nextMonth}
class="w-8 h-8 flex items-center justify-center hover:bg-gray-800 rounded-md text-gray-300 touch-manipulation">
<ChevronRight size={15} /></button>
</div>
</div>

<div class="grid grid-cols-7 gap-px sm:gap-1">
{#each dayNames as dn, i}
<div class="text-center py-1.5">
<span class="text-[10px] font-semibold text-gray-500 sm:hidden">{dayNamesShort[i]}</span>
<span class="hidden sm:inline text-xs font-semibold text-gray-400">{dn}</span>
</div>
{/each}
{#each calendarDays as day}
{#if day === null}
<div class="aspect-square"></div>
{:else}
{@const dt      = getTasksForDate(day)}
{@const isToday = isDateToday(day)}
{@const isSel   = selectedDate === day}
{@const isPast  = isDatePast(day)}
{@const hasOD   = dt.some(t => isOverdue(t))}
<button onclick={() => selectedDate = day}
class="aspect-square rounded sm:rounded-md transition-all touch-manipulation flex flex-col items-center justify-start pt-1 sm:pt-1.5 px-0.5
{isToday ? 'bg-blue-600 ring-1 ring-blue-400' : ''}
{isSel && !isToday ? 'bg-blue-950 ring-1 ring-blue-500' : ''}
{!isToday && !isSel ? 'hover:bg-gray-800' : ''}
{hasOD && !isToday ? 'ring-1 ring-red-700/60' : ''}">
<span class="text-[11px] sm:text-sm font-medium leading-none
{isToday ? 'text-white' : isSel ? 'text-blue-300' : isPast ? 'text-gray-500' : 'text-gray-300'}">
{day}
</span>
{#if dt.length > 0}
<div class="flex gap-px mt-0.5 sm:mt-1 flex-wrap justify-center">
{#each dt.slice(0,4) as task}
<div class="w-1 h-1 sm:w-1.5 sm:h-1.5 rounded-full {getDotColor(task, isToday)}"></div>
{/each}
{#if dt.length > 4}<span class="text-[8px] {isToday?'text-white':'text-gray-500'}">+</span>{/if}
</div>
{/if}
</button>
{/if}
{/each}
</div>

<!-- Legend -->
<div class="mt-2 pt-2 border-t border-gray-800 flex items-center gap-3 flex-wrap">
<div class="flex items-center gap-1"><div class="w-2 h-2 rounded-full bg-gray-500"></div><span class="text-[10px] text-gray-500">Task</span></div>
<div class="flex items-center gap-1"><div class="w-2 h-2 rounded-full bg-blue-400"></div><span class="text-[10px] text-gray-500">Event</span></div>
<div class="flex items-center gap-1"><div class="w-2 h-2 rounded-full bg-purple-400"></div><span class="text-[10px] text-gray-500">Reminder</span></div>
<div class="flex items-center gap-1"><div class="w-2 h-2 rounded-full bg-amber-400"></div><span class="text-[10px] text-gray-500">Payment Due</span></div>
<div class="flex items-center gap-1"><div class="w-2 h-2 rounded-full bg-green-400"></div><span class="text-[10px] text-gray-500">Paid</span></div>
<div class="flex items-center gap-1"><div class="w-2 h-2 rounded-full bg-red-500"></div><span class="text-[10px] text-gray-500">Overdue</span></div>
</div>
</div>

<!-- Sidebar -->
<div class="sm:col-span-1 bg-gray-900 border border-gray-800 rounded-lg p-2.5 sm:p-3 flex flex-col gap-2">

<!-- Month Summary Stats -->
{#if tasks.filter(t => { if(!t.due_date) return false; const d=new Date(t.due_date); return d.getMonth()===currentMonth && d.getFullYear()===currentYear; }).length > 0}
{@const ms = tasks.filter(t => { if(!t.due_date) return false; const d=new Date(t.due_date); return d.getMonth()===currentMonth && d.getFullYear()===currentYear; })}
<div class="bg-gray-800 rounded-md p-2 border border-gray-700">
<p class="text-[10px] font-semibold text-gray-400 mb-1.5 flex items-center gap-1">
<CalendarIcon size={10} class="text-blue-400" /> {monthNames[currentMonth]} Summary
</p>
<div class="grid grid-cols-4 gap-1">
<div class="flex flex-col items-center bg-gray-900 rounded p-1.5 border border-gray-700">
<span class="text-base font-bold text-white leading-none">{ms.length}</span>
<span class="text-[9px] text-gray-500 mt-0.5">Total</span>
</div>
<div class="flex flex-col items-center bg-gray-900 rounded p-1.5 border border-gray-700">
<span class="text-base font-bold text-gray-400 leading-none">{ms.filter(t=>t.status==='todo').length}</span>
<span class="text-[9px] text-gray-500 mt-0.5">Todo</span>
</div>
<div class="flex flex-col items-center bg-gray-900 rounded p-1.5 border border-gray-700">
<span class="text-base font-bold text-blue-400 leading-none">{ms.filter(t=>t.status==='in-progress').length}</span>
<span class="text-[9px] text-gray-500 mt-0.5">Active</span>
</div>
<div class="flex flex-col items-center bg-gray-900 rounded p-1.5 border border-gray-700">
<span class="text-base font-bold text-green-400 leading-none">{ms.filter(t=>t.status==='done').length}</span>
<span class="text-[9px] text-gray-500 mt-0.5">Done</span>
</div>
</div>
</div>
{/if}

<!-- Payment Summary -->
{#if paymentStats.count > 0}
{@const pct = paymentStats.totalAmt > 0 ? Math.round((paymentStats.paidAmt / paymentStats.totalAmt) * 100) : 0}
<div class="bg-gray-800 rounded-md p-2 border border-amber-900/40">
<p class="text-[10px] font-semibold text-amber-400 mb-2 flex items-center gap-1">
<CreditCard size={10} /> Payments — {monthNames[currentMonth]}
</p>
<!-- Progress bar -->
<div class="h-1.5 bg-gray-700 rounded-full mb-2 overflow-hidden">
<div class="h-full bg-green-500 rounded-full transition-all" style="width:{pct}%"></div>
</div>
<div class="space-y-1">
<div class="flex justify-between text-[10px]">
<span class="text-gray-400">Total ({paymentStats.count})</span>
<span class="text-white font-mono">{formatCurrency(paymentStats.totalAmt, paymentStats.currency)}</span>
</div>
<div class="flex justify-between text-[10px]">
<span class="text-green-400">Paid ({paymentStats.paidCount})</span>
<span class="text-green-400 font-mono">{formatCurrency(paymentStats.paidAmt, paymentStats.currency)}</span>
</div>
<div class="flex justify-between text-[10px]">
<span class="text-amber-400">Outstanding</span>
<span class="text-amber-400 font-mono">{formatCurrency(paymentStats.unpaidAmt, paymentStats.currency)}</span>
</div>
{#if paymentStats.overdueCount > 0}
<div class="flex items-center gap-1 mt-1 bg-red-950/40 rounded px-1.5 py-1 border border-red-900/40">
<AlertCircle size={10} class="text-red-400 flex-shrink-0" />
<span class="text-[10px] text-red-300">{paymentStats.overdueCount} overdue</span>
</div>
{/if}
</div>
</div>
{/if}

<!-- Date task list -->
<div class="flex-1 min-h-0">
<div class="flex items-center justify-between mb-2">
<h3 class="text-xs font-bold text-white leading-none">
{#if selectedDate}{monthNames[currentMonth]} {selectedDate}{:else}Select a date{/if}
</h3>
{#if selectedDate}
<button onclick={() => openAddTaskModal(selectedDate)}
class="h-7 px-2 bg-blue-600 hover:bg-blue-700 text-white rounded text-[10px] font-medium flex items-center gap-1 touch-manipulation">
<Plus size={10} /> Add
</button>
{/if}
</div>

{#if selectedDate}
{#if selectedDateTasks.length > 0}
<div class="space-y-1.5 overflow-y-auto max-h-[50vh] sm:max-h-[calc(100vh-420px)] pr-0.5">
{#each selectedDateTasks as task}
{@const sub = (task.calendar_subtype || 'regular') as CalendarSub}
{@const SIcon = subIcon(sub)}
{@const overdue = isOverdue(task)}

{#if editingTask && editingTask.id === task.id}
<div class="border border-blue-500 rounded-lg p-2.5 bg-gray-800 space-y-2">
<input type="text" bind:value={editingTask.title}
class="w-full h-9 px-3 bg-gray-900 border border-gray-700 text-white rounded-md text-sm focus:ring-1 focus:ring-blue-500" placeholder="Title" />
<textarea bind:value={editingTask.description}
class="w-full px-3 py-2 bg-gray-900 border border-gray-700 text-white rounded-md text-sm focus:ring-1 focus:ring-blue-500" rows="2" placeholder="Description"></textarea>
<div class="grid grid-cols-2 gap-1.5">
<select bind:value={editingTask.calendar_subtype}
class="h-9 px-2 bg-gray-900 border border-gray-700 text-white rounded-md text-xs focus:ring-1 focus:ring-blue-500">
<option value="regular">Regular</option>
<option value="payment">Payment</option>
<option value="reminder">Reminder</option>
<option value="event">Event</option>
</select>
<select bind:value={editingTask.priority}
class="h-9 px-2 bg-gray-900 border border-gray-700 text-white rounded-md text-xs focus:ring-1 focus:ring-blue-500">
<option value="low">Low</option>
<option value="medium">Medium</option>
<option value="high">High</option>
</select>
</div>
{#if editingTask.calendar_subtype === 'payment' || editingTask.is_payment}
<div class="grid grid-cols-2 gap-1.5">
<input type="number" bind:value={editingTask.amount} min="0" step="0.01"
class="h-9 px-2 bg-gray-900 border border-gray-700 text-white rounded-md text-xs focus:ring-1 focus:ring-amber-500" placeholder="Amount" />
<select bind:value={editingTask.currency}
class="h-9 px-2 bg-gray-900 border border-gray-700 text-white rounded-md text-xs focus:ring-1 focus:ring-amber-500">
{#each currencies as c}<option value={c}>{c}</option>{/each}
</select>
</div>
<label class="flex items-center gap-2 text-xs text-gray-300 cursor-pointer">
<input type="checkbox" bind:checked={editingTask.is_paid} class="rounded accent-green-500" />
Mark as Paid
</label>
{/if}
<input type="date" bind:value={editingTask.due_date}
class="w-full h-9 px-3 bg-gray-900 border border-gray-700 text-white rounded-md text-sm focus:ring-1 focus:ring-blue-500" />
<div class="flex gap-1.5">
<button onclick={handleSaveEdit} class="flex-1 h-9 bg-blue-600 text-white rounded-md text-xs font-medium hover:bg-blue-700 touch-manipulation">Save</button>
<button onclick={() => { editingTask = null; }} class="flex-1 h-9 bg-gray-700 text-gray-200 rounded-md text-xs font-medium hover:bg-gray-600 touch-manipulation">Cancel</button>
</div>
</div>
{:else}
<div class="border-l-2 {getPriorityBorder(task.priority)} rounded-md p-2 flex items-start gap-2
{overdue ? 'ring-1 ring-red-700/40' : ''}">
<SIcon size={14} class="{subIconColor(task)} mt-0.5 flex-shrink-0" />
<div class="flex-1 min-w-0">
<p class="text-xs font-semibold truncate {(task.status === 'done' || task.is_paid) ? 'line-through text-gray-500' : 'text-gray-200'}">{task.title}</p>
{#if task.description}
<p class="text-[10px] text-gray-400 mt-0.5 line-clamp-1">{task.description}</p>
{/if}
{#if task.is_payment || sub === 'payment'}
<div class="flex items-center gap-1.5 mt-1 flex-wrap">
{#if task.amount && task.amount > 0}
<span class="text-[10px] font-mono font-semibold {task.is_paid ? 'text-green-400' : overdue ? 'text-red-400' : 'text-amber-400'}">
{formatCurrency(task.amount, task.currency || 'USD')}
</span>
{/if}
{#if task.is_paid}
<span class="text-[9px] px-1.5 py-0.5 rounded-full bg-green-950 text-green-300 border border-green-800">Paid</span>
{#if task.paid_at}<span class="text-[9px] text-gray-500">{formatDate(task.paid_at)}</span>{/if}
{:else if overdue}
<span class="text-[9px] px-1.5 py-0.5 rounded-full bg-red-950 text-red-300 border border-red-800 flex items-center gap-0.5">
<AlertCircle size={8} /> Overdue
</span>
{:else}
<span class="text-[9px] px-1.5 py-0.5 rounded-full bg-amber-950 text-amber-300 border border-amber-800">Unpaid</span>
{/if}
</div>
{/if}
<div class="flex items-center gap-1.5 mt-1 flex-wrap">
{#if sub !== 'regular'}
<span class="text-[9px] px-1.5 py-0.5 rounded border {getSubtypeBadge(sub)}">{sub}</span>
{/if}
{#if task.priority === 'high'}
<span class="text-[9px] px-1.5 py-0.5 rounded-full bg-red-950 text-red-300">High</span>
{/if}
</div>
</div>
<div class="flex flex-col gap-0.5 flex-shrink-0">
{#if (sub === 'payment' || task.is_payment) && !task.is_paid}
<button onclick={() => showMarkPaid = task}
class="w-7 h-7 flex items-center justify-center hover:bg-green-950 rounded text-green-400 touch-manipulation" title="Mark Paid">
<Check size={13} />
</button>
{/if}
{#if (sub === 'payment' || task.is_payment) && task.is_paid}
<button onclick={() => markAsUnpaid(task)}
class="w-7 h-7 flex items-center justify-center hover:bg-yellow-950 rounded text-yellow-400 touch-manipulation" title="Undo Paid">
<RotateCcw size={11} />
</button>
{/if}
<button onclick={() => { editingTask = { ...task }; }}
class="w-7 h-7 flex items-center justify-center hover:bg-blue-950 rounded text-blue-400 touch-manipulation" title="Edit">
<Edit2 size={12} />
</button>
<button onclick={() => handleArchive(task.id)}
class="w-7 h-7 flex items-center justify-center hover:bg-gray-700 rounded text-gray-400 touch-manipulation" title="Archive">
<Archive size={12} />
</button>
<button onclick={() => handleDelete(task.id)}
class="w-7 h-7 flex items-center justify-center hover:bg-red-950 rounded text-red-400 touch-manipulation" title="Delete">
<Trash2 size={12} />
</button>
</div>
</div>
{/if}
{/each}
</div>
{:else}
<div class="text-center py-6">
<CalendarIcon size={28} class="text-gray-700 mx-auto mb-2" />
<p class="text-xs text-gray-500">Nothing here</p>
<button onclick={() => openAddTaskModal(selectedDate)}
class="inline-flex items-center gap-1 mt-2 h-8 px-3 bg-blue-600 text-white rounded-md text-xs font-medium hover:bg-blue-700 touch-manipulation">
<Plus size={12} /> Add
</button>
</div>
{/if}
{:else}
<div class="text-center py-6">
<CalendarIcon size={28} class="text-gray-700 mx-auto mb-2" />
<p class="text-xs text-gray-500">Tap a date to view</p>
</div>
{/if}
</div>
</div>
</div>

<!-- ════════════════════ PAYMENT HISTORY TAB ════════════════════════════ -->
{:else if activeTab === 'history'}
<div class="bg-gray-900 border border-gray-800 rounded-lg p-3 sm:p-4">
<div class="flex flex-wrap items-center justify-between gap-2 mb-4">
<div>
<h2 class="text-sm font-bold text-white flex items-center gap-2">
<CreditCard size={16} class="text-amber-400" /> Payment History
</h2>
<p class="text-[11px] text-gray-500 mt-0.5">All payment tasks — active & archived</p>
</div>
<div class="flex items-center gap-2 flex-wrap">
<div class="flex gap-1 bg-gray-800 rounded-lg p-0.5">
{#each [['all','All'],['paid','Paid'],['unpaid','Unpaid'],['overdue','Overdue']] as [val, label]}
<button onclick={() => historyFilter = val as HistoryFilter}
class="px-2.5 py-1 rounded-md text-[11px] font-medium transition
{historyFilter === val ? 'bg-blue-600 text-white' : 'text-gray-400 hover:text-white'}">
{label}
</button>
{/each}
</div>
<button onclick={() => downloadCSV(allPaymentTasks, `payment-history-${new Date().toISOString().slice(0,10)}.csv`)}
class="flex items-center gap-1.5 h-8 px-3 bg-green-700 hover:bg-green-600 text-white rounded-lg text-xs font-medium transition touch-manipulation">
<Download size={12} /> Export CSV
</button>
</div>
</div>

{#if allPaymentTasks.length > 0}
<!-- All-time stats -->
{@const all = [...tasks, ...archivedTasks].filter(t => t.is_payment || t.calendar_subtype === 'payment')}
{@const totAmt  = all.reduce((s,t)=>s+(t.amount||0),0)}
{@const paidAmt = all.filter(t=>t.is_paid).reduce((s,t)=>s+(t.amount||0),0)}
<div class="grid grid-cols-2 sm:grid-cols-4 gap-2 mb-5">
<div class="bg-gray-800 rounded-lg p-3 border border-gray-700">
<p class="text-[10px] text-gray-500 mb-1">Total Records</p>
<p class="text-2xl font-bold text-white">{all.length}</p>
</div>
<div class="bg-gray-800 rounded-lg p-3 border border-gray-700">
<p class="text-[10px] text-gray-500 mb-1">Total Value</p>
<p class="text-lg font-bold text-white font-mono">{formatCurrency(totAmt, all[0]?.currency||'USD')}</p>
</div>
<div class="bg-green-950/40 rounded-lg p-3 border border-green-900/40">
<p class="text-[10px] text-green-500 mb-1">Paid ({all.filter(t=>t.is_paid).length})</p>
<p class="text-lg font-bold text-green-400 font-mono">{formatCurrency(paidAmt, all[0]?.currency||'USD')}</p>
</div>
<div class="bg-amber-950/30 rounded-lg p-3 border border-amber-900/40">
<p class="text-[10px] text-amber-500 mb-1">Outstanding ({all.filter(t=>!t.is_paid).length})</p>
<p class="text-lg font-bold text-amber-400 font-mono">{formatCurrency(totAmt-paidAmt, all[0]?.currency||'USD')}</p>
</div>
</div>

<!-- Grouped by month -->
{#each [...groupByMonth(allPaymentTasks)] as [month, monthTasks]}
{@const mTot  = monthTasks.reduce((s,t)=>s+(t.amount||0),0)}
{@const mPaid = monthTasks.filter(t=>t.is_paid).reduce((s,t)=>s+(t.amount||0),0)}
<div class="mb-5">
<div class="flex items-center gap-2 mb-2">
<h3 class="text-xs font-bold text-gray-300 flex-shrink-0">{month}</h3>
<div class="h-px flex-1 bg-gray-800"></div>
<span class="text-[10px] text-gray-500 flex-shrink-0">
{formatCurrency(mPaid, monthTasks[0]?.currency||'USD')} / {formatCurrency(mTot, monthTasks[0]?.currency||'USD')} paid
</span>
</div>
<div class="space-y-1.5">
{#each monthTasks as task}
{@const od = isOverdue(task)}
<div class="flex items-center gap-3 bg-gray-800 border rounded-lg px-3 py-2.5 transition
{od ? 'border-red-800/50' : task.is_paid ? 'border-green-900/30' : 'border-gray-700'}">
<div class="flex-shrink-0">
{#if task.is_paid}
<div class="w-7 h-7 rounded-full bg-green-950 border border-green-700 flex items-center justify-center">
<Check size={13} class="text-green-400" />
</div>
{:else if od}
<div class="w-7 h-7 rounded-full bg-red-950 border border-red-700 flex items-center justify-center">
<AlertCircle size={13} class="text-red-400" />
</div>
{:else}
<div class="w-7 h-7 rounded-full bg-amber-950 border border-amber-700 flex items-center justify-center">
<CreditCard size={13} class="text-amber-400" />
</div>
{/if}
</div>
<div class="flex-1 min-w-0">
<p class="text-sm font-semibold truncate {task.is_paid ? 'text-gray-400 line-through' : 'text-white'}">{task.title}</p>
<div class="flex items-center gap-2 flex-wrap mt-0.5">
<span class="text-[10px] text-gray-500">Due: {formatDate(task.due_date)}</span>
{#if task.is_paid && task.paid_at}
<span class="text-[10px] text-green-500">Paid: {formatDate(task.paid_at)}</span>
{/if}
{#if task.payment_notes}
<span class="text-[10px] text-gray-500 italic">"{task.payment_notes}"</span>
{/if}
{#if task.is_archived}
<span class="text-[9px] px-1.5 py-0.5 rounded bg-gray-700 text-gray-400">archived</span>
{/if}
</div>
</div>
<div class="flex items-center gap-2 flex-shrink-0">
<div class="text-right">
<p class="text-sm font-bold font-mono {task.is_paid ? 'text-green-400' : od ? 'text-red-400' : 'text-amber-400'}">
{task.amount ? formatCurrency(task.amount, task.currency||'USD') : '—'}
</p>
{#if task.is_paid}
<span class="text-[9px] text-green-500">PAID</span>
{:else if od}
<span class="text-[9px] text-red-500">OVERDUE</span>
{:else}
<span class="text-[9px] text-amber-500">PENDING</span>
{/if}
</div>
{#if !task.is_paid}
<button onclick={() => showMarkPaid = task}
class="h-7 px-2.5 bg-green-700 hover:bg-green-600 text-white rounded text-[11px] font-medium touch-manipulation">
Mark Paid
</button>
{:else}
<button onclick={() => markAsUnpaid(task)}
class="h-7 px-2 bg-gray-700 hover:bg-gray-600 text-gray-300 rounded text-[10px] touch-manipulation">
Undo
</button>
{/if}
</div>
</div>
{/each}
</div>
</div>
{/each}
{:else}
<div class="text-center py-16">
<CreditCard size={40} class="text-gray-700 mx-auto mb-3" />
<p class="text-sm text-gray-500">No payment tasks yet.</p>
<p class="text-xs text-gray-600 mt-1">Create a task with type "Payment" in the calendar.</p>
</div>
{/if}
</div>

<!-- ════════════════════ ARCHIVE TAB ════════════════════════════════════ -->
{:else if activeTab === 'archive'}
<div class="bg-gray-900 border border-gray-800 rounded-lg p-3 sm:p-4">
<div class="flex flex-wrap items-center justify-between gap-2 mb-4">
<div>
<h2 class="text-sm font-bold text-white flex items-center gap-2">
<Archive size={16} class="text-gray-400" /> Calendar Archive
</h2>
<p class="text-[11px] text-gray-500 mt-0.5">Archived events, reminders, tasks & payment history</p>
</div>
<div class="flex items-center gap-2 flex-wrap">
<input bind:value={archiveSearch} type="text" placeholder="Search..."
class="h-8 px-3 bg-gray-800 border border-gray-700 text-white rounded-lg text-xs w-36 focus:ring-1 focus:ring-blue-500" />
<div class="flex gap-1 bg-gray-800 rounded-lg p-0.5">
{#each [['all','All'],['payment','Payments'],['event','Events'],['reminder','Reminders']] as [f, label]}
<button onclick={() => archiveFilter = f as typeof archiveFilter}
class="px-2 py-1 rounded-md text-[11px] font-medium transition
{archiveFilter === f ? 'bg-gray-600 text-white' : 'text-gray-400 hover:text-white'}">
{label}
</button>
{/each}
</div>
<button onclick={() => downloadCSV(filteredArchive, `archive-${new Date().toISOString().slice(0,10)}.csv`)}
class="flex items-center gap-1.5 h-8 px-3 bg-gray-700 hover:bg-gray-600 text-white rounded-lg text-xs font-medium transition touch-manipulation">
<Download size={12} /> CSV
</button>
</div>
</div>

{#if filteredArchive.length > 0}
<div class="space-y-1.5">
{#each filteredArchive as task}
{@const sub = (task.calendar_subtype || 'regular') as CalendarSub}
{@const SIcon = subIcon(sub)}
<div class="flex items-center gap-3 bg-gray-800/60 border border-gray-700/60 rounded-lg px-3 py-2.5 hover:bg-gray-800 transition">
<SIcon size={15} class="{subIconColor(task)} flex-shrink-0" />
<div class="flex-1 min-w-0">
<div class="flex items-center gap-2 flex-wrap">
<p class="text-xs font-semibold text-gray-300 truncate">{task.title}</p>
{#if sub !== 'regular'}
<span class="text-[9px] px-1.5 py-0.5 rounded border {getSubtypeBadge(sub)}">{sub}</span>
{/if}
{#if (task.is_payment || sub === 'payment') && task.amount}
<span class="text-[10px] font-mono {task.is_paid ? 'text-green-400' : 'text-amber-400'}">{formatCurrency(task.amount, task.currency||'USD')}</span>
{#if task.is_paid}
<span class="text-[9px] px-1.5 py-0.5 rounded bg-green-950 text-green-400 border border-green-800">Paid</span>
{/if}
{/if}
</div>
<div class="flex items-center gap-2 mt-0.5 flex-wrap">
{#if task.due_date}<span class="text-[10px] text-gray-500">Due: {formatDate(task.due_date)}</span>{/if}
{#if task.paid_at}<span class="text-[10px] text-green-600">Paid: {formatDate(task.paid_at)}</span>{/if}
{#if task.description}<span class="text-[10px] text-gray-600 truncate max-w-[220px]">{task.description}</span>{/if}
</div>
</div>
<div class="flex items-center gap-1 flex-shrink-0">
<button onclick={() => handleRestore(task.id)}
class="h-7 px-2 bg-blue-950 hover:bg-blue-900 text-blue-300 rounded text-[10px] font-medium flex items-center gap-1 touch-manipulation">
<RotateCcw size={10} /> Restore
</button>
<button onclick={() => handlePermDelete(task.id)}
class="w-7 h-7 flex items-center justify-center hover:bg-red-950 rounded text-red-500 touch-manipulation" title="Delete forever">
<Trash2 size={12} />
</button>
</div>
</div>
{/each}
</div>
{:else}
<div class="text-center py-16">
<Archive size={40} class="text-gray-700 mx-auto mb-3" />
<p class="text-sm text-gray-500">
{#if archiveSearch || archiveFilter !== 'all'}No results match your filter.{:else}Archive is empty.{/if}
</p>
</div>
{/if}
</div>
{/if}

{/if}

<!-- ════════ Add Task Modal ════════════════════════════════════════════════ -->
{#if showAddTask}
<div class="fixed inset-0 bg-black/60 z-50 flex items-end sm:items-center justify-center sm:p-4"
onclick={(e) => { if (e.target === e.currentTarget) closeAddTaskModal(); }}>
<div class="bg-gray-900 border border-gray-800 rounded-t-2xl sm:rounded-2xl shadow-2xl w-full sm:max-w-md p-4 sm:p-5 max-h-[92vh] overflow-y-auto">
<div class="flex items-center justify-between mb-3">
<h3 class="text-sm font-bold text-white">Add to Calendar</h3>
<button onclick={closeAddTaskModal} class="w-8 h-8 flex items-center justify-center hover:bg-gray-800 rounded-lg text-gray-400"><X size={18} /></button>
</div>
<div class="space-y-3">
<!-- Type selector -->
<div>
<label class="block text-xs font-medium text-gray-400 mb-1.5">Type</label>
<div class="grid grid-cols-4 gap-1 bg-gray-800 rounded-lg p-1">
{#each [['regular','Task',Circle],['payment','Payment',CreditCard],['reminder','Reminder',Bell],['event','Event',Star]] as [val, lbl, Ico]}
<button type="button"
onclick={() => { newTask.calendar_subtype = val as CalendarSub; newTask.is_payment = (val==='payment'); }}
class="flex flex-col items-center gap-1 py-2 rounded-md text-[10px] font-medium transition
{newTask.calendar_subtype === val
? (val==='payment' ? 'bg-amber-700 text-white' : val==='reminder' ? 'bg-purple-700 text-white' : val==='event' ? 'bg-blue-700 text-white' : 'bg-gray-600 text-white')
: 'text-gray-400 hover:text-white'}">
<Ico size={14} />{lbl}
</button>
{/each}
</div>
</div>

<div>
<label for="t-title" class="block text-xs font-medium text-gray-400 mb-1">Title *</label>
<input id="t-title" type="text" bind:value={newTask.title}
class="w-full h-10 px-3 bg-gray-800 border border-gray-700 text-white rounded-lg text-sm focus:ring-2 focus:ring-blue-500"
placeholder="{newTask.calendar_subtype==='payment'?'e.g. Rent, Netflix, Insurance…':newTask.calendar_subtype==='reminder'?'Reminder title…':'Task title'}" />
</div>

<div>
<label for="t-desc" class="block text-xs font-medium text-gray-400 mb-1">Description</label>
<textarea id="t-desc" bind:value={newTask.description}
class="w-full px-3 py-2 bg-gray-800 border border-gray-700 text-white rounded-lg text-sm focus:ring-2 focus:ring-blue-500" rows="2" placeholder="Optional notes"></textarea>
</div>

<!-- Payment fields -->
{#if newTask.calendar_subtype === 'payment'}
<div class="bg-amber-950/20 border border-amber-900/40 rounded-lg p-3 space-y-2.5">
<p class="text-[10px] font-semibold text-amber-400 flex items-center gap-1.5"><CreditCard size={10} /> Payment Details</p>
<div class="grid grid-cols-2 gap-2">
<div>
<label class="block text-[10px] font-medium text-gray-500 mb-1">Amount</label>
<input type="number" bind:value={newTask.amount} min="0" step="0.01"
class="w-full h-9 px-2 bg-gray-800 border border-gray-700 text-white rounded-md text-sm focus:ring-1 focus:ring-amber-500" placeholder="0.00" />
</div>
<div>
<label class="block text-[10px] font-medium text-gray-500 mb-1">Currency</label>
<select bind:value={newTask.currency}
class="w-full h-9 px-2 bg-gray-800 border border-gray-700 text-white rounded-md text-sm focus:ring-1 focus:ring-amber-500">
{#each currencies as c}<option value={c}>{c}</option>{/each}
</select>
</div>
</div>
<div>
<label class="block text-[10px] font-medium text-gray-500 mb-1">Payment Notes</label>
<input type="text" bind:value={newTask.payment_notes}
class="w-full h-9 px-2 bg-gray-800 border border-gray-700 text-white rounded-md text-sm focus:ring-1 focus:ring-amber-500" placeholder="e.g. Bank transfer, account #…" />
</div>
<label class="flex items-center gap-2 cursor-pointer">
<input type="checkbox" bind:checked={newTask.is_paid} class="rounded accent-green-500" />
<span class="text-xs text-gray-300">Already paid</span>
</label>
</div>
{/if}

<div class="grid grid-cols-2 gap-2">
<div>
<label for="t-priority" class="block text-xs font-medium text-gray-400 mb-1">Priority</label>
<select id="t-priority" bind:value={newTask.priority}
class="w-full h-10 px-3 bg-gray-800 border border-gray-700 text-white rounded-lg text-sm focus:ring-2 focus:ring-blue-500">
<option value="low">Low</option><option value="medium">Medium</option><option value="high">High</option>
</select>
</div>
<div>
<label for="t-status" class="block text-xs font-medium text-gray-400 mb-1">Status</label>
<select id="t-status" bind:value={newTask.status}
class="w-full h-10 px-3 bg-gray-800 border border-gray-700 text-white rounded-lg text-sm focus:ring-2 focus:ring-blue-500">
<option value="todo">To Do</option><option value="in-progress">In Progress</option><option value="done">Done</option>
</select>
</div>
</div>

<div>
<label for="t-date" class="block text-xs font-medium text-gray-400 mb-1">
{newTask.calendar_subtype === 'payment' ? 'Payment Due Date' : 'Date'}
</label>
<input id="t-date" type="date" bind:value={newTask.due_date}
class="w-full h-10 px-3 bg-gray-800 border border-gray-700 text-white rounded-lg text-sm focus:ring-2 focus:ring-blue-500" />
</div>

<!-- Recurring -->
<div class="bg-gray-800 border border-gray-700 rounded-lg p-3 space-y-2.5">
<label class="flex items-center gap-2.5 cursor-pointer select-none">
<div class="relative flex-shrink-0">
<input type="checkbox" bind:checked={recurring} class="sr-only peer" />
<div class="w-9 h-5 rounded-full transition-colors {recurring ? 'bg-blue-600' : 'bg-gray-600'}"></div>
<div class="absolute top-0.5 left-0.5 w-4 h-4 bg-white rounded-full shadow transition-transform {recurring ? 'translate-x-4' : 'translate-x-0'}"></div>
</div>
<span class="text-xs font-medium text-gray-300">Recurring</span>
</label>
{#if recurring}
<div class="grid grid-cols-2 gap-2">
<div>
<label class="block text-[10px] font-medium text-gray-500 mb-1">Repeat</label>
<select bind:value={recurringType}
class="w-full h-9 px-2 bg-gray-700 border border-gray-600 text-white rounded-md text-xs">
<option value="daily">Daily</option><option value="weekly">Weekly</option>
<option value="monthly">Monthly</option><option value="yearly">Yearly</option>
</select>
</div>
<div>
<label class="block text-[10px] font-medium text-gray-500 mb-1">Occurrences</label>
<input type="number" bind:value={recurringCount} min="2" max="60"
class="w-full h-9 px-2 bg-gray-700 border border-gray-600 text-white rounded-md text-xs" />
</div>
</div>
<p class="text-[10px] text-blue-400">Creates {recurringCount} × {recurringType} starting from the chosen date.</p>
{/if}
</div>

<div class="flex gap-2 pt-1">
<button onclick={handleAddTask}
class="flex-1 h-10 {newTask.calendar_subtype==='payment'?'bg-amber-600 hover:bg-amber-700':newTask.calendar_subtype==='reminder'?'bg-purple-600 hover:bg-purple-700':newTask.calendar_subtype==='event'?'bg-blue-600 hover:bg-blue-700':'bg-gray-600 hover:bg-gray-700'} text-white rounded-lg text-sm font-medium active:scale-95 transition touch-manipulation">
{recurring ? `Create ${recurringCount}` : 'Create'} {newTask.calendar_subtype === 'payment' ? 'Payment' : newTask.calendar_subtype === 'reminder' ? 'Reminder' : newTask.calendar_subtype === 'event' ? 'Event' : 'Task'}
</button>
<button onclick={closeAddTaskModal}
class="flex-1 h-10 bg-gray-700 text-gray-200 rounded-lg text-sm font-medium hover:bg-gray-600 touch-manipulation">Cancel</button>
</div>
</div>
</div>
</div>
{/if}

<!-- ════════ Mark Paid Confirmation Modal ══════════════════════════════════ -->
{#if showMarkPaid}
{@const t = showMarkPaid}
<div class="fixed inset-0 bg-black/70 z-50 flex items-center justify-center p-4">
<div class="bg-gray-900 border border-green-800 rounded-2xl shadow-2xl w-full max-w-sm p-5">
<div class="flex items-center gap-3 mb-4">
<div class="w-10 h-10 rounded-full bg-green-950 border border-green-700 flex items-center justify-center flex-shrink-0">
<Check size={20} class="text-green-400" />
</div>
<div>
<h3 class="text-sm font-bold text-white">Confirm Payment</h3>
<p class="text-xs text-gray-400 mt-0.5">This will mark the item as paid</p>
</div>
</div>
<div class="bg-gray-800 rounded-lg p-3 mb-4">
<p class="text-sm font-semibold text-white">{t.title}</p>
{#if t.amount}
<p class="text-xl font-bold text-green-400 font-mono mt-1">{formatCurrency(t.amount, t.currency||'USD')}</p>
{/if}
{#if t.due_date}
<p class="text-xs text-gray-500 mt-0.5">Due: {formatDate(t.due_date)}</p>
{/if}
{#if t.payment_notes}
<p class="text-xs text-gray-500 mt-0.5 italic">"{t.payment_notes}"</p>
{/if}
</div>
<p class="text-xs text-gray-400 mb-4">
Marking paid on: <span class="text-white font-medium">
{new Date().toLocaleDateString('en-US',{year:'numeric',month:'long',day:'numeric'})}
</span>
</p>
<div class="flex gap-2">
<button onclick={() => markAsPaid(t)}
class="flex-1 h-10 bg-green-600 hover:bg-green-700 text-white rounded-lg text-sm font-bold transition touch-manipulation">
Confirm Paid ✓
</button>
<button onclick={() => showMarkPaid = null}
class="flex-1 h-10 bg-gray-700 hover:bg-gray-600 text-gray-200 rounded-lg text-sm font-medium touch-manipulation">
Cancel
</button>
</div>
</div>
</div>
{/if}
</div>
