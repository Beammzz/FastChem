"use client";

import { useState } from "react";
import { AdminFindingsResponse, AdminRule } from "@/types";
import { BarList, Empty, Panel, Pill, Spinner, Table, formatDateTime, formatNumber } from "./ui";

interface AnticheatPanelProps {
  rules: AdminRule[];
  findings: AdminFindingsResponse | null;
  loading: boolean;
  filters: { rule: string; mode: string; action: string };
  page: number;
  onFilter: (filters: { rule: string; mode: string; action: string }) => void;
  onPage: (page: number) => void;
  onSaveRule: (rule: {
    name: string;
    enabled: boolean;
    action: string;
    params: Record<string, number>;
  }) => Promise<void>;
}

const MODES = ["", "casual", "match", "ranked", "room"];
const ACTIONS = ["", "observe", "reject"];

export default function AnticheatPanel({
  rules,
  findings,
  loading,
  filters,
  page,
  onFilter,
  onPage,
  onSaveRule,
}: AnticheatPanelProps) {
  const pages = findings ? Math.max(1, Math.ceil(findings.total / findings.pageSize)) : 1;

  return (
    <div className="space-y-4">
      <Panel
        title="กฎตรวจจับ"
        subtitle="แก้แล้วมีผลทันที ไม่ต้องรีสตาร์ต — observe คือบันทึกอย่างเดียว, reject คือตัดคะแนนข้อนั้น"
      >
        <div className="space-y-3">
          {rules.map((rule) => (
            <RuleEditor key={rule.name} rule={rule} onSave={onSaveRule} />
          ))}
          {rules.length === 0 && <Empty icon="🛡️" text="ยังไม่มีกฎ" />}
        </div>
      </Panel>

      {findings && findings.dropped > 0 && (
        <div className="bg-amber-500/10 border border-amber-500/20 rounded-xl px-4 py-3 text-xs text-amber-300">
          ทิ้ง findings ไป {formatNumber(findings.dropped)} รายการ เพราะบัฟเฟอร์เต็ม —
          ตัวเลขในตารางนี้ต่ำกว่าความเป็นจริง
        </div>
      )}

      {findings && (
        <div className="grid sm:grid-cols-3 gap-4">
          <Panel title="ตามกฎ">
            <BarList items={findings.byRule.map((b) => ({ label: b.label, count: b.count }))} />
          </Panel>
          <Panel title="ตามโหมด">
            <BarList items={findings.byMode.map((b) => ({ label: b.label, count: b.count }))} />
          </Panel>
          <Panel title="ตามการจัดการ">
            <BarList items={findings.byAction.map((b) => ({ label: b.label, count: b.count }))} />
          </Panel>
        </div>
      )}

      <Panel
        title="Findings"
        subtitle={findings ? `${formatNumber(findings.total)} รายการที่ตรงเงื่อนไข` : undefined}
        action={
          <div className="flex flex-wrap gap-1.5">
            <Select
              value={filters.rule}
              options={["", ...rules.map((r) => r.name)]}
              placeholder="ทุกกฎ"
              onChange={(rule) => onFilter({ ...filters, rule })}
            />
            <Select
              value={filters.mode}
              options={MODES}
              placeholder="ทุกโหมด"
              onChange={(mode) => onFilter({ ...filters, mode })}
            />
            <Select
              value={filters.action}
              options={ACTIONS}
              placeholder="ทุกการจัดการ"
              onChange={(action) => onFilter({ ...filters, action })}
            />
          </div>
        }
      >
        {loading ? (
          <Spinner />
        ) : !findings || findings.findings.length === 0 ? (
          <Empty icon="✅" text="ยังไม่มี finding ที่ตรงเงื่อนไข" />
        ) : (
          <>
            <Table
              head={
                <>
                  <th className="text-left px-4 sm:px-5 py-3">เวลา</th>
                  <th className="text-left px-3 py-3">ผู้เล่น</th>
                  <th className="text-left px-3 py-3">กฎ</th>
                  <th className="text-left px-3 py-3">รายละเอียด</th>
                  <th className="text-left px-3 py-3">โหมด</th>
                  <th className="text-right px-4 sm:px-5 py-3">เวลาที่ใช้</th>
                </>
              }
            >
              {findings.findings.map((f) => (
                <tr key={f.id} className="border-b border-white/5 last:border-0">
                  <td className="px-4 sm:px-5 py-2.5 text-xs text-gray-500 whitespace-nowrap">
                    {formatDateTime(f.at)}
                  </td>
                  <td className="px-3 py-2.5">
                    {f.username ? (
                      <span className="text-gray-200">{f.username}</span>
                    ) : (
                      <span className="text-gray-500 text-xs" title={f.subject}>
                        ไม่ระบุตัวตน
                      </span>
                    )}
                  </td>
                  <td className="px-3 py-2.5">
                    <Pill tone={f.action === "reject" ? "red" : "amber"}>{f.rule}</Pill>
                  </td>
                  <td className="px-3 py-2.5 text-xs text-gray-400 max-w-[280px]">{f.detail}</td>
                  <td className="px-3 py-2.5 text-xs text-gray-500">
                    {f.mode}
                    {f.matchId !== 0 && <span className="text-gray-600"> #{f.matchId}</span>}
                  </td>
                  <td className="px-4 sm:px-5 py-2.5 text-right tabular-nums text-gray-400">
                    {f.timeSpent.toFixed(2)} วิ
                  </td>
                </tr>
              ))}
            </Table>

            <div className="flex items-center justify-between mt-4 text-xs text-gray-500">
              <span>
                หน้า {page} / {pages}
              </span>
              <div className="flex gap-2">
                <button
                  disabled={page <= 1}
                  onClick={() => onPage(page - 1)}
                  className="px-3 py-1 rounded-md border border-white/10 disabled:opacity-30 hover:text-white cursor-pointer disabled:cursor-default"
                >
                  ก่อนหน้า
                </button>
                <button
                  disabled={page >= pages}
                  onClick={() => onPage(page + 1)}
                  className="px-3 py-1 rounded-md border border-white/10 disabled:opacity-30 hover:text-white cursor-pointer disabled:cursor-default"
                >
                  ถัดไป
                </button>
              </div>
            </div>
          </>
        )}
      </Panel>
    </div>
  );
}

function Select({
  value,
  options,
  placeholder,
  onChange,
}: {
  value: string;
  options: string[];
  placeholder: string;
  onChange: (value: string) => void;
}) {
  return (
    <select
      value={value}
      onChange={(e) => onChange(e.target.value)}
      className="bg-[#0a0a1a] border border-white/10 rounded-lg px-2 py-1.5 text-xs text-gray-300 focus:outline-none focus:border-violet-500/40 cursor-pointer"
    >
      {options.map((option) => (
        <option key={option} value={option}>
          {option === "" ? placeholder : option}
        </option>
      ))}
    </select>
  );
}

function RuleEditor({
  rule,
  onSave,
}: {
  rule: AdminRule;
  onSave: (rule: {
    name: string;
    enabled: boolean;
    action: string;
    params: Record<string, number>;
  }) => Promise<void>;
}) {
  const [enabled, setEnabled] = useState(rule.enabled);
  const [action, setAction] = useState(rule.action);
  const [params, setParams] = useState<Record<string, string>>(
    Object.fromEntries(Object.entries(rule.params ?? {}).map(([k, v]) => [k, String(v)]))
  );
  const [error, setError] = useState("");
  const [busy, setBusy] = useState(false);

  const dirty =
    enabled !== rule.enabled ||
    action !== rule.action ||
    Object.entries(params).some(([k, v]) => Number(v) !== rule.params[k]);

  const save = async () => {
    setBusy(true);
    setError("");
    try {
      const numeric: Record<string, number> = {};
      for (const [key, value] of Object.entries(params)) {
        const parsed = Number(value);
        if (Number.isFinite(parsed)) numeric[key] = parsed;
      }
      await onSave({ name: rule.name, enabled, action, params: numeric });
    } catch (e) {
      setError(e instanceof Error ? e.message : "บันทึกไม่สำเร็จ");
    } finally {
      setBusy(false);
    }
  };

  return (
    <div className="border border-white/5 rounded-xl p-4 bg-[#0a0a1a]/40">
      <div className="flex flex-wrap items-center justify-between gap-3">
        <div>
          <div className="flex items-center gap-2">
            <span className="font-medium text-sm">{rule.name}</span>
            {action === "reject" && <Pill tone="red">ตัดคะแนนจริง</Pill>}
            {!enabled && <Pill tone="muted">ปิดอยู่</Pill>}
          </div>
          <div className="text-xs text-gray-500 mt-0.5">
            จับได้ {formatNumber(rule.findings)} ครั้ง · 24 ชม. {formatNumber(rule.findings24h)}
          </div>
        </div>

        <div className="flex items-center gap-2">
          <label className="flex items-center gap-1.5 text-xs text-gray-400 cursor-pointer">
            <input
              type="checkbox"
              checked={enabled}
              onChange={(e) => setEnabled(e.target.checked)}
              className="accent-violet-500 cursor-pointer"
            />
            เปิดใช้
          </label>
          <select
            value={action}
            onChange={(e) => setAction(e.target.value)}
            className="bg-[#0a0a1a] border border-white/10 rounded-lg px-2 py-1.5 text-xs focus:outline-none focus:border-violet-500/40 cursor-pointer"
          >
            <option value="observe">observe</option>
            <option value="reject">reject</option>
          </select>
          <button
            disabled={!dirty || busy}
            onClick={save}
            className="px-3 py-1.5 rounded-lg text-xs font-medium bg-violet-500/20 border border-violet-500/30 text-violet-300 hover:bg-violet-500/30 cursor-pointer disabled:opacity-30 disabled:cursor-default"
          >
            บันทึก
          </button>
        </div>
      </div>

      {Object.keys(params).length > 0 && (
        <div className="grid grid-cols-2 sm:grid-cols-4 gap-2 mt-3">
          {Object.entries(params).map(([key, value]) => (
            <label key={key} className="block">
              <span className="text-[11px] text-gray-500">{key}</span>
              <input
                type="number"
                step="any"
                value={value}
                onChange={(e) => setParams({ ...params, [key]: e.target.value })}
                className="mt-0.5 w-full bg-[#0a0a1a] border border-white/10 rounded-lg px-2 py-1 text-xs tabular-nums focus:outline-none focus:border-violet-500/40"
              />
            </label>
          ))}
        </div>
      )}

      {action === "reject" && enabled && (
        <p className="text-[11px] text-amber-400/80 mt-3">
          กฎนี้จะทำให้คำตอบที่เข้าเงื่อนไขถูกนับว่าผิดและไม่ได้คะแนน
          ผู้เล่นจะไม่ถูกแจ้งว่าโดนกฎข้อไหน
        </p>
      )}

      {error && <p className="text-xs text-red-400 mt-2">{error}</p>}
    </div>
  );
}
