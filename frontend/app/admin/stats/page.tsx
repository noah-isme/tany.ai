const popularQuestions = [
  { question: "Berapa harga pembuatan website?", hits: 126 },
  { question: "Berapa lama pengerjaan AI assistant?", hits: 93 },
  { question: "Apakah ada paket maintenance?", hits: 54 },
];

const metrics = [
  { label: "Traffic mingguan", value: "3.240", delta: "+12%" },
  { label: "Lead baru", value: "38", delta: "+5" },
  { label: "Rasio respons", value: "98%", delta: "+2%" },
];

const sparkline = [45, 52, 39, 60, 72, 64, 88];

export default function StatsPage() {
  return (
    <div className="space-y-6">
      <div className="space-y-2">
        <h1 className="text-xl font-semibold text-slate-900 dark:text-slate-100">Statistik</h1>
        <p className="text-sm text-slate-600 dark:text-slate-400">
          Placeholder insight sesuai blueprint: pantau tren pertanyaan, traffic, dan lead sebelum integrasi analitik penuh.
        </p>
      </div>

      <section className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
        {metrics.map((metric) => (
          <div
            key={metric.label}
            className="rounded-2xl border border-slate-200 bg-white/80 p-4 shadow-sm dark:border-slate-800 dark:bg-slate-900/70"
          >
            <p className="text-xs uppercase tracking-[0.35em] text-slate-400">{metric.label}</p>
            <p className="mt-3 text-2xl font-semibold text-slate-900 dark:text-slate-100">{metric.value}</p>
            <p className="mt-1 text-xs text-emerald-400">{metric.delta} dibanding minggu lalu</p>
          </div>
        ))}
      </section>

      <section className="grid gap-4 lg:grid-cols-[2fr_1fr]">
        <div className="rounded-2xl border border-slate-200 bg-white/80 p-6 shadow-sm dark:border-slate-800 dark:bg-slate-900/70">
          <h2 className="text-base font-semibold text-slate-900 dark:text-slate-100">Pertanyaan Terpopuler</h2>
          <ul className="mt-4 space-y-3">
            {popularQuestions.map((item) => (
              <li key={item.question} className="flex items-center justify-between">
                <span className="text-sm text-slate-600 dark:text-slate-300">{item.question}</span>
                <span className="text-sm font-semibold text-indigo-400">{item.hits}x</span>
              </li>
            ))}
          </ul>
        </div>
        <div className="rounded-2xl border border-slate-200 bg-white/80 p-6 shadow-sm dark:border-slate-800 dark:bg-slate-900/70">
          <h2 className="text-base font-semibold text-slate-900 dark:text-slate-100">Traffic 7 Hari</h2>
          <div className="mt-4 flex h-32 items-end justify-between gap-1">
            {sparkline.map((value, index) => (
              <div
                key={index}
                className="w-full rounded-t-md bg-indigo-500/70"
                style={{ height: `${value}%` }}
                aria-hidden="true"
              />
            ))}
          </div>
          <p className="mt-3 text-xs text-slate-500 dark:text-slate-400">
            Visual placeholder: akan digantikan chart interaktif setelah integrasi analitik backend.
          </p>
        </div>
      </section>
    </div>
  );
}
