"use client";

import {
  Activity,
  ClipboardList,
  FileText,
  Loader2,
  Play,
  RefreshCcw,
  Send,
  Sparkles,
} from "lucide-react";
import { FormEvent, useEffect, useMemo, useState } from "react";
import {
  ApiEnvelope,
  HealthStatus,
  ScriptResponse,
  VideoTask,
  requestJSON,
} from "@/lib/api";

type TaskListPayload = {
  items: VideoTask[];
};

const t = {
  appName: "\u0041\u0049 \u6570\u5b57\u4eba",
  appSubtitle: "\u89c6\u9891\u751f\u4ea7\u540e\u53f0",
  navTasks: "\u89c6\u9891\u4efb\u52a1",
  navScript: "\u6587\u6848\u751f\u6210",
  navHealth: "\u670d\u52a1\u72b6\u6001",
  title: "\u8fd0\u8425\u63a7\u5236\u53f0",
  phase: "\u7b2c\u4e00\u9636\u6bb5\u0020\u004d\u0056\u0050",
  productKind: "\u6570\u5b57\u4eba\u53e3\u64ad\u77ed\u89c6\u9891",
  refresh: "\u5237\u65b0",
  scriptDraft: "\u6587\u6848\u8349\u7a3f",
  generate: "\u751f\u6210\u6587\u6848",
  product: "\u5546\u54c1/\u4e3b\u9898",
  points: "\u5356\u70b9",
  tone: "\u8bed\u6c14",
  duration: "\u65f6\u957f\uff08\u79d2\uff09",
  createTask: "\u521b\u5efa\u4efb\u52a1",
  queue: "\u52a0\u5165\u961f\u5217",
  avatar: "\u6570\u5b57\u4eba\u5934\u50cf URL",
  voiceover: "\u53e3\u64ad\u6587\u6848",
  totalPrefix: "\u5171",
  totalSuffix: "\u6761",
  status: "\u72b6\u6001",
  createdAt: "\u521b\u5efa\u65f6\u95f4",
  action: "\u64cd\u4f5c",
  noTasks: "\u6682\u65e0\u4efb\u52a1",
  openTask: "\u67e5\u770b\u4efb\u52a1",
  taskDetail: "\u4efb\u52a1\u8be6\u60c5",
  taskID: "\u4efb\u52a1 ID",
  script: "\u6587\u6848",
  latestDraft: "\u6700\u65b0\u6587\u6848\u8349\u7a3f",
  noDraft: "\u8fd8\u6ca1\u6709\u751f\u6210\u6587\u6848\u8349\u7a3f",
  noSelected: "\u5c1a\u672a\u9009\u62e9\u4efb\u52a1",
  resultVideo: "\u6210\u679c\u89c6\u9891",
  openResult: "\u6253\u5f00\u6210\u679c",
  ossObject: "OSS \u5bf9\u8c61",
  previewDuration: "\u9884\u89c8\u65f6\u957f",
  renderCommand: "\u6e32\u67d3\u547d\u4ee4",
  errorMessage: "\u9519\u8bef\u4fe1\u606f",
  seconds: "\u79d2",
};

const messages = {
  healthFailed: "\u670d\u52a1\u72b6\u6001\u68c0\u67e5\u5931\u8d25",
  tasksFailed: "\u4efb\u52a1\u5217\u8868\u5237\u65b0\u5931\u8d25",
  scriptGenerated: "\u6587\u6848\u8349\u7a3f\u5df2\u751f\u6210",
  scriptFailed: "\u6587\u6848\u751f\u6210\u5931\u8d25",
  taskQueued: "\u89c6\u9891\u4efb\u52a1\u5df2\u52a0\u5165\u961f\u5217",
  taskFailed: "\u521b\u5efa\u4efb\u52a1\u5931\u8d25",
};

const defaultScript =
  "\u4ecb\u7ecd\u6d3b\u52a8\u4eae\u70b9\uff0c\u8bf4\u660e\u6838\u5fc3\u4ef7\u503c\uff0c\u6700\u540e\u7ed9\u51fa\u660e\u786e\u7684\u884c\u52a8\u5f15\u5bfc\u3002";

export default function Dashboard() {
  const [health, setHealth] = useState<HealthStatus[]>([]);
  const [tasks, setTasks] = useState<VideoTask[]>([]);
  const [selectedTaskID, setSelectedTaskID] = useState<string>("");
  const [productName, setProductName] = useState("\u5468\u672b\u5496\u5561\u56e2\u8d2d\u5238");
  const [sellingPoints, setSellingPoints] = useState("\u4e0b\u5355\u65b9\u4fbf, \u5230\u5e97\u5373\u7528, \u5468\u672b\u9650\u65f6\u4f18\u60e0");
  const [tone, setTone] = useState("\u6e29\u6696\u81ea\u7136");
  const [script, setScript] = useState(defaultScript);
  const [avatarImageURL, setAvatarImageURL] = useState("https://example.com/avatar.png");
  const [durationSeconds, setDurationSeconds] = useState(60);
  const [scriptDraft, setScriptDraft] = useState<ScriptResponse | null>(null);
  const [message, setMessage] = useState("");
  const [loadingHealth, setLoadingHealth] = useState(false);
  const [loadingTasks, setLoadingTasks] = useState(false);
  const [generatingScript, setGeneratingScript] = useState(false);
  const [creatingTask, setCreatingTask] = useState(false);

  const selectedTask = useMemo(
    () => tasks.find((task) => task.id === selectedTaskID) ?? tasks[0],
    [selectedTaskID, tasks],
  );

  useEffect(() => {
    void refreshAll();
  }, []);

  async function refreshAll() {
    await Promise.all([loadHealth(), loadTasks()]);
  }

  async function loadHealth() {
    setLoadingHealth(true);
    try {
      const data = await requestJSON<{ services: HealthStatus[] }>("/api/health");
      setHealth(data.services);
    } catch (error) {
      setMessage(error instanceof Error ? error.message : messages.healthFailed);
    } finally {
      setLoadingHealth(false);
    }
  }

  async function loadTasks() {
    setLoadingTasks(true);
    try {
      const envelope = await requestJSON<ApiEnvelope<TaskListPayload>>("/api/backend/video");
      setTasks(envelope.data.items);
      if (!selectedTaskID && envelope.data.items.length > 0) {
        setSelectedTaskID(envelope.data.items[0].id);
      }
    } catch (error) {
      setMessage(error instanceof Error ? error.message : messages.tasksFailed);
    } finally {
      setLoadingTasks(false);
    }
  }

  async function generateScript(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setGeneratingScript(true);
    setMessage("");
    try {
      const response = await requestJSON<ScriptResponse>("/api/ai/script", {
        method: "POST",
        body: JSON.stringify({
          product_name: productName,
          selling_points: sellingPoints
            .split(",")
            .map((point) => point.trim())
            .filter(Boolean),
          tone,
          duration_seconds: durationSeconds,
        }),
      });
      setScriptDraft(response);
      setScript(response.scenes.map((scene) => scene.voiceover).join(" "));
      setMessage(messages.scriptGenerated);
    } catch (error) {
      setMessage(error instanceof Error ? error.message : messages.scriptFailed);
    } finally {
      setGeneratingScript(false);
    }
  }

  async function createTask(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setCreatingTask(true);
    setMessage("");
    try {
      const envelope = await requestJSON<ApiEnvelope<VideoTask>>("/api/backend/video/create", {
        method: "POST",
        body: JSON.stringify({
          product_name: productName,
          script,
          avatar_image_url: avatarImageURL,
          duration_seconds: durationSeconds,
        }),
      });
      setSelectedTaskID(envelope.data.id);
      setMessage(messages.taskQueued);
      await loadTasks();
    } catch (error) {
      setMessage(error instanceof Error ? error.message : messages.taskFailed);
    } finally {
      setCreatingTask(false);
    }
  }

  return (
    <main className="min-h-screen bg-canvas text-ink">
      <div className="mx-auto flex w-full max-w-[1440px]">
        <aside className="hidden min-h-screen w-64 border-r border-line bg-panel px-5 py-6 lg:block">
          <div className="mb-8 flex items-center gap-3">
            <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-teal-50 text-teal-700">
              <Sparkles size={21} aria-hidden="true" />
            </div>
            <div>
              <div className="text-base font-semibold leading-tight">{t.appName}</div>
              <div className="text-xs font-medium text-muted">{t.appSubtitle}</div>
            </div>
          </div>

          <nav className="space-y-1 text-sm font-medium text-muted">
            <NavLink href="#tasks" active icon={<ClipboardList size={16} aria-hidden="true" />} label={t.navTasks} />
            <NavLink href="#script" icon={<FileText size={16} aria-hidden="true" />} label={t.navScript} />
            <NavLink href="#health" icon={<Activity size={16} aria-hidden="true" />} label={t.navHealth} />
          </nav>
        </aside>

        <section className="min-w-0 flex-1 px-4 py-5 sm:px-6 lg:px-8">
          <header className="mb-5 flex flex-col gap-4 border-b border-line pb-5 xl:flex-row xl:items-center xl:justify-between">
            <div>
              <h1 className="text-2xl font-semibold leading-tight tracking-normal sm:text-3xl">{t.title}</h1>
              <div className="mt-2 flex flex-wrap gap-2 text-sm text-muted">
                <span>{t.phase}</span>
                <span>/</span>
                <span>{t.productKind}</span>
              </div>
            </div>
            <button
              type="button"
              onClick={() => void refreshAll()}
              className="inline-flex h-10 items-center justify-center gap-2 rounded-md border border-line bg-panel px-4 text-sm font-semibold shadow-panel hover:bg-canvas disabled:cursor-not-allowed disabled:opacity-60"
              disabled={loadingHealth || loadingTasks}
            >
              {loadingHealth || loadingTasks ? (
                <Loader2 className="animate-spin" size={16} aria-hidden="true" />
              ) : (
                <RefreshCcw size={16} aria-hidden="true" />
              )}
              {t.refresh}
            </button>
          </header>

          {message ? (
            <div className="mb-5 rounded-md border border-line bg-panel px-4 py-3 text-sm font-medium text-ink shadow-panel">
              {message}
            </div>
          ) : null}

          <section id="health" className="mb-5 grid gap-3 md:grid-cols-3">
            {health.map((service) => (
              <div key={service.name} className="rounded-md border border-line bg-panel p-4 shadow-panel">
                <div className="flex items-center justify-between gap-3">
                  <div>
                    <h2 className="text-sm font-semibold">{service.name}</h2>
                    <p className="mt-1 break-all text-xs text-muted">{service.url}</p>
                  </div>
                  <span
                    className={[
                      "inline-flex items-center rounded-md px-2 py-1 text-xs font-bold",
                      service.ok ? "bg-teal-50 text-teal-700" : "bg-amber-50 text-amber-600",
                    ].join(" ")}
                  >
                    {service.detail}
                  </span>
                </div>
              </div>
            ))}
          </section>

          <div className="grid gap-5 xl:grid-cols-[minmax(0,440px)_minmax(0,1fr)]">
            <section id="script" className="space-y-5">
              <form onSubmit={generateScript} className="rounded-md border border-line bg-panel p-4 shadow-panel">
                <div className="mb-4 flex items-center justify-between gap-3">
                  <h2 className="text-base font-semibold">{t.scriptDraft}</h2>
                  <ActionButton loading={generatingScript} icon={<Sparkles size={16} aria-hidden="true" />}>
                    {t.generate}
                  </ActionButton>
                </div>

                <Field label={t.product}>
                  <input value={productName} onChange={(event) => setProductName(event.target.value)} className="h-10 w-full rounded-md border border-line bg-white px-3 text-sm" />
                </Field>

                <Field label={t.points}>
                  <textarea value={sellingPoints} onChange={(event) => setSellingPoints(event.target.value)} rows={3} className="w-full resize-none rounded-md border border-line bg-white px-3 py-2 text-sm" />
                </Field>

                <div className="grid gap-3 sm:grid-cols-2">
                  <Field label={t.tone}>
                    <input value={tone} onChange={(event) => setTone(event.target.value)} className="h-10 w-full rounded-md border border-line bg-white px-3 text-sm" />
                  </Field>
                  <Field label={t.duration}>
                    <input type="number" min={15} max={180} value={durationSeconds} onChange={(event) => setDurationSeconds(Number(event.target.value))} className="h-10 w-full rounded-md border border-line bg-white px-3 text-sm" />
                  </Field>
                </div>
              </form>

              <form onSubmit={createTask} className="rounded-md border border-line bg-panel p-4 shadow-panel">
                <div className="mb-4 flex items-center justify-between gap-3">
                  <h2 className="text-base font-semibold">{t.createTask}</h2>
                  <button type="submit" className="inline-flex h-9 items-center justify-center gap-2 rounded-md bg-ink px-3 text-sm font-semibold text-white hover:bg-[#2a3140] disabled:cursor-not-allowed disabled:opacity-60" disabled={creatingTask}>
                    {creatingTask ? <Loader2 className="animate-spin" size={16} aria-hidden="true" /> : <Send size={16} aria-hidden="true" />}
                    {t.queue}
                  </button>
                </div>

                <Field label={t.avatar}>
                  <input value={avatarImageURL} onChange={(event) => setAvatarImageURL(event.target.value)} className="h-10 w-full rounded-md border border-line bg-white px-3 text-sm" />
                </Field>

                <Field label={t.voiceover}>
                  <textarea value={script} onChange={(event) => setScript(event.target.value)} rows={6} className="w-full resize-none rounded-md border border-line bg-white px-3 py-2 text-sm" />
                </Field>
              </form>
            </section>

            <section className="space-y-5">
              <section id="tasks" className="rounded-md border border-line bg-panel shadow-panel">
                <div className="flex items-center justify-between gap-3 border-b border-line px-4 py-3">
                  <h2 className="text-base font-semibold">{t.navTasks}</h2>
                  <span className="text-sm font-medium text-muted">{t.totalPrefix} {tasks.length} {t.totalSuffix}</span>
                </div>
                <div className="overflow-x-auto">
                  <table className="w-full min-w-[760px] border-collapse text-sm">
                    <thead className="bg-canvas text-left text-xs text-muted">
                      <tr>
                        <th className="px-4 py-3 font-semibold">{t.product}</th>
                        <th className="px-4 py-3 font-semibold">{t.status}</th>
                        <th className="px-4 py-3 font-semibold">{t.duration}</th>
                        <th className="px-4 py-3 font-semibold">{t.createdAt}</th>
                        <th className="px-4 py-3 font-semibold">{t.action}</th>
                      </tr>
                    </thead>
                    <tbody>
                      {tasks.map((task) => (
                        <tr key={task.id} className="border-t border-line">
                          <td className="px-4 py-3 font-medium">{task.product_name}</td>
                          <td className="px-4 py-3">
                            <span className={statusClass(task.status)}>{formatStatus(task.status)}</span>
                          </td>
                          <td className="px-4 py-3 text-muted">{task.duration_seconds} {t.seconds}</td>
                          <td className="px-4 py-3 text-muted">{formatDate(task.created_at)}</td>
                          <td className="px-4 py-3">
                            <button type="button" onClick={() => setSelectedTaskID(task.id)} className="inline-flex h-8 w-8 items-center justify-center rounded-md border border-line hover:bg-canvas" title={t.openTask} aria-label={t.openTask}>
                              <Play size={15} aria-hidden="true" />
                            </button>
                          </td>
                        </tr>
                      ))}
                      {tasks.length === 0 ? (
                        <tr>
                          <td className="px-4 py-8 text-center text-muted" colSpan={5}>{t.noTasks}</td>
                        </tr>
                      ) : null}
                    </tbody>
                  </table>
                </div>
              </section>

              <section className="grid gap-5 lg:grid-cols-2">
                <div className="rounded-md border border-line bg-panel p-4 shadow-panel">
                  <h2 className="mb-3 text-base font-semibold">{t.taskDetail}</h2>
                  {selectedTask ? (
                    <div className="space-y-3 text-sm">
                      <Detail label={t.taskID} value={selectedTask.id} />
                      <Detail label={t.product} value={selectedTask.product_name} />
                      <Detail label={t.avatar} value={selectedTask.avatar_image_url} />
                      <Detail label={t.script} value={selectedTask.script} />
                      {selectedTask.output_url ? (
                        <div>
                          <div className="mb-1 text-xs font-bold text-muted">{t.resultVideo}</div>
                          <a className="inline-flex rounded-md bg-teal-700 px-3 py-2 text-sm font-semibold text-white hover:bg-teal-600" href={selectedTask.output_url} target="_blank" rel="noreferrer">
                            {t.openResult}
                          </a>
                        </div>
                      ) : null}
                      {selectedTask.object_key ? <Detail label={t.ossObject} value={selectedTask.object_key} /> : null}
                      {selectedTask.preview_seconds ? <Detail label={t.previewDuration} value={`${selectedTask.preview_seconds} ${t.seconds}`} /> : null}
                      {selectedTask.error_message ? <Detail label={t.errorMessage} value={selectedTask.error_message} /> : null}
                      {selectedTask.render_command?.length ? <Detail label={t.renderCommand} value={selectedTask.render_command.join(" ")} /> : null}
                    </div>
                  ) : (
                    <p className="text-sm text-muted">{t.noSelected}</p>
                  )}
                </div>

                <div className="rounded-md border border-line bg-panel p-4 shadow-panel">
                  <h2 className="mb-3 text-base font-semibold">{t.latestDraft}</h2>
                  {scriptDraft ? (
                    <div className="space-y-3">
                      {scriptDraft.scenes.map((scene) => (
                        <div key={scene.order} className="border-b border-line pb-3 last:border-b-0 last:pb-0">
                          <div className="mb-1 flex items-center justify-between gap-3 text-sm font-semibold">
                            <span>{scene.title}</span>
                            <span className="text-xs text-muted">{scene.duration_seconds} {t.seconds}</span>
                          </div>
                          <p className="text-sm leading-6 text-muted">{scene.voiceover}</p>
                        </div>
                      ))}
                    </div>
                  ) : (
                    <p className="text-sm text-muted">{t.noDraft}</p>
                  )}
                </div>
              </section>
            </section>
          </div>
        </section>
      </div>
    </main>
  );
}

function NavLink({ href, icon, label, active = false }: { href: string; icon: React.ReactNode; label: string; active?: boolean }) {
  return (
    <a className={["flex items-center gap-3 rounded-md px-3 py-2", active ? "bg-teal-50 text-teal-700" : "hover:bg-canvas"].join(" ")} href={href}>
      {icon}
      {label}
    </a>
  );
}

function ActionButton({ loading, icon, children }: { loading: boolean; icon: React.ReactNode; children: React.ReactNode }) {
  return (
    <button type="submit" className="inline-flex h-9 items-center justify-center gap-2 rounded-md bg-teal-700 px-3 text-sm font-semibold text-white hover:bg-teal-600 disabled:cursor-not-allowed disabled:opacity-60" disabled={loading}>
      {loading ? <Loader2 className="animate-spin" size={16} aria-hidden="true" /> : icon}
      {children}
    </button>
  );
}

function Field({ label, children }: { label: string; children: React.ReactNode }) {
  return (
    <label className="mb-3 block">
      <span className="mb-1.5 block text-sm font-semibold text-ink">{label}</span>
      {children}
    </label>
  );
}

function Detail({ label, value }: { label: string; value: string }) {
  return (
    <div>
      <div className="mb-1 text-xs font-bold text-muted">{label}</div>
      <div className="break-words rounded-md bg-canvas px-3 py-2 leading-6">{value}</div>
    </div>
  );
}

function formatDate(value: string) {
  return new Intl.DateTimeFormat("zh-CN", {
    month: "2-digit",
    day: "2-digit",
    hour: "2-digit",
    minute: "2-digit",
  }).format(new Date(value));
}

function formatStatus(status: VideoTask["status"]) {
  const labels: Record<VideoTask["status"], string> = {
    queued: "\u6392\u961f\u4e2d",
    running: "\u5904\u7406\u4e2d",
    succeeded: "\u5df2\u5b8c\u6210",
    failed: "\u5931\u8d25",
  };
  return labels[status] ?? status;
}

function statusClass(status: VideoTask["status"]) {
  const base = "rounded-md px-2 py-1 text-xs font-bold";
  if (status === "succeeded") return `${base} bg-teal-50 text-teal-700`;
  if (status === "running") return `${base} bg-blue-50 text-blue-700`;
  if (status === "failed") return `${base} bg-amber-50 text-amber-600`;
  return `${base} bg-slate-100 text-slate-700`;
}
