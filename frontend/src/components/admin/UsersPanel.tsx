"use client";

import { useState } from "react";
import Link from "next/link";
import { AdminUserRow, AdminUsersResponse, AdminUserUpdate } from "@/types";
import { Empty, Panel, Pill, Spinner, Table, formatDateTime, formatNumber } from "./ui";

interface UsersPanelProps {
  data: AdminUsersResponse | null;
  loading: boolean;
  search: string;
  sort: string;
  page: number;
  onSearch: (value: string) => void;
  onSort: (key: string) => void;
  onPage: (page: number) => void;
  onUpdate: (update: AdminUserUpdate) => Promise<void>;
  onDelete: (userId: number) => Promise<void>;
  currentUserId?: number;
}

const COLUMNS: { key: string; label: string; align: "left" | "right" }[] = [
  { key: "username", label: "ผู้เล่น", align: "left" },
  { key: "points", label: "คะแนน", align: "right" },
  { key: "rating", label: "Rating", align: "right" },
  { key: "games", label: "เกม", align: "right" },
  { key: "accuracy", label: "แม่นยำ", align: "right" },
  { key: "findings", label: "Flags", align: "right" },
  { key: "lastplay", label: "เล่นล่าสุด", align: "right" },
];

export default function UsersPanel({
  data,
  loading,
  search,
  sort,
  page,
  onSearch,
  onSort,
  onPage,
  onUpdate,
  onDelete,
  currentUserId,
}: UsersPanelProps) {
  const [editing, setEditing] = useState<AdminUserRow | null>(null);

  const pages = data ? Math.max(1, Math.ceil(data.total / data.pageSize)) : 1;

  return (
    <div className="space-y-4">
      <Panel
        title="ผู้เล่น"
        subtitle={data ? `${formatNumber(data.total)} บัญชี` : undefined}
        action={
          <input
            value={search}
            onChange={(e) => onSearch(e.target.value)}
            placeholder="ค้นหาชื่อผู้เล่น"
            className="bg-[#0a0a1a] border border-white/10 rounded-lg px-3 py-1.5 text-xs w-40 sm:w-56 focus:outline-none focus:border-violet-500/40"
          />
        }
      >
        {loading ? (
          <Spinner />
        ) : !data || data.users.length === 0 ? (
          <Empty icon="👤" text="ไม่พบผู้เล่น" />
        ) : (
          <>
            <Table
              head={
                <>
                  {COLUMNS.map((col) => (
                    <th
                      key={col.key}
                      onClick={() => onSort(col.key)}
                      className={`px-3 py-3 cursor-pointer select-none hover:text-gray-300 ${
                        col.align === "right" ? "text-right" : "text-left"
                      } ${sort === col.key ? "text-violet-400" : ""}`}
                    >
                      {col.label}
                      {sort === col.key && " ↓"}
                    </th>
                  ))}
                  <th className="px-3 py-3 text-right">จัดการ</th>
                </>
              }
            >
              {data.users.map((u) => (
                <tr key={u.id} className="border-b border-white/5 last:border-0 hover:bg-white/[0.02]">
                  <td className="px-3 py-2.5">
                    <div className="flex items-center gap-2">
                      {u.online && (
                        <span className="w-1.5 h-1.5 rounded-full bg-emerald-400" title="กำลังออนไลน์" />
                      )}
                      <Link
                        href={`/profile/${u.username}/`}
                        className="font-medium text-gray-100 hover:text-violet-400 truncate max-w-[140px]"
                      >
                        {u.username}
                      </Link>
                      {u.isAdmin && <Pill tone="violet">admin</Pill>}
                    </div>
                    <div className="text-[11px] text-gray-600 mt-0.5">
                      สมัคร {formatDateTime(u.createdAt)}
                    </div>
                  </td>
                  <td className="px-3 py-2.5 text-right tabular-nums">{formatNumber(u.totalPoints)}</td>
                  <td className="px-3 py-2.5 text-right tabular-nums">
                    {u.rating}
                    <span className="text-[11px] text-gray-600 ml-1">
                      {u.rankedWins}W/{u.rankedLosses}L
                    </span>
                  </td>
                  <td className="px-3 py-2.5 text-right tabular-nums text-gray-400">{u.games}</td>
                  <td className="px-3 py-2.5 text-right tabular-nums text-gray-400">
                    {u.totalAnswered > 0 ? `${u.accuracy.toFixed(0)}%` : "—"}
                  </td>
                  <td className="px-3 py-2.5 text-right tabular-nums">
                    {u.findings > 0 ? (
                      <span className="text-amber-400">{u.findings}</span>
                    ) : (
                      <span className="text-gray-600">0</span>
                    )}
                  </td>
                  <td className="px-3 py-2.5 text-right text-xs text-gray-500">
                    {formatDateTime(u.lastPlayed)}
                  </td>
                  <td className="px-3 py-2.5 text-right">
                    <button
                      onClick={() => setEditing(u)}
                      className="text-xs text-violet-400 hover:text-violet-300 cursor-pointer"
                    >
                      แก้ไข
                    </button>
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

      {editing && (
        <EditUserDialog
          user={editing}
          isSelf={editing.id === currentUserId}
          onClose={() => setEditing(null)}
          onUpdate={onUpdate}
          onDelete={onDelete}
        />
      )}
    </div>
  );
}

interface EditUserDialogProps {
  user: AdminUserRow;
  isSelf: boolean;
  onClose: () => void;
  onUpdate: (update: AdminUserUpdate) => Promise<void>;
  onDelete: (userId: number) => Promise<void>;
}

function EditUserDialog({ user, isSelf, onClose, onUpdate, onDelete }: EditUserDialogProps) {
  const [rating, setRating] = useState(String(user.rating));
  const [points, setPoints] = useState(String(user.totalPoints));
  const [password, setPassword] = useState("");
  const [confirmDelete, setConfirmDelete] = useState("");
  const [error, setError] = useState("");
  const [busy, setBusy] = useState(false);

  // Only changed fields are sent — the endpoint leaves omitted ones alone, so a
  // half-filled form cannot blank a value the operator never touched.
  const run = async (action: () => Promise<void>) => {
    setBusy(true);
    setError("");
    try {
      await action();
      onClose();
    } catch (e) {
      setError(e instanceof Error ? e.message : "ทำรายการไม่สำเร็จ");
    } finally {
      setBusy(false);
    }
  };

  const save = () =>
    run(async () => {
      const update: AdminUserUpdate = { userId: user.id };
      if (Number(rating) !== user.rating) update.rating = Number(rating);
      if (Number(points) !== user.totalPoints) update.totalPoints = Number(points);
      if (password) update.password = password;
      if (Object.keys(update).length === 1) return;
      await onUpdate(update);
    });

  return (
    <div className="fixed inset-0 z-50 bg-black/70 backdrop-blur-sm flex items-center justify-center p-4">
      <div className="bg-[#12122a] border border-white/10 rounded-2xl w-full max-w-md max-h-[90vh] overflow-y-auto">
        <header className="px-5 py-4 border-b border-white/5 flex items-center justify-between">
          <h3 className="font-semibold">{user.username}</h3>
          <button onClick={onClose} className="text-gray-500 hover:text-white cursor-pointer">
            ✕
          </button>
        </header>

        <div className="p-5 space-y-4 text-sm">
          <label className="block">
            <span className="text-xs text-gray-500">Rating</span>
            <input
              type="number"
              value={rating}
              onChange={(e) => setRating(e.target.value)}
              className="mt-1 w-full bg-[#0a0a1a] border border-white/10 rounded-lg px-3 py-2 tabular-nums focus:outline-none focus:border-violet-500/40"
            />
          </label>

          <label className="block">
            <span className="text-xs text-gray-500">คะแนนรวม</span>
            <input
              type="number"
              value={points}
              onChange={(e) => setPoints(e.target.value)}
              className="mt-1 w-full bg-[#0a0a1a] border border-white/10 rounded-lg px-3 py-2 tabular-nums focus:outline-none focus:border-violet-500/40"
            />
          </label>

          <label className="block">
            <span className="text-xs text-gray-500">ตั้งรหัสผ่านใหม่ (อย่างน้อย 6 ตัว)</span>
            <input
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              placeholder="เว้นว่างไว้ถ้าไม่เปลี่ยน"
              className="mt-1 w-full bg-[#0a0a1a] border border-white/10 rounded-lg px-3 py-2 focus:outline-none focus:border-violet-500/40"
            />
          </label>

          <div className="flex items-center justify-between border-t border-white/5 pt-4">
            <div>
              <div className="text-gray-300">สิทธิ์แอดมิน</div>
              <div className="text-xs text-gray-500">
                {isSelf ? "ถอดสิทธิ์ตัวเองไม่ได้" : "เข้าถึงคอนโซลนี้ได้ทั้งหมด"}
              </div>
            </div>
            <button
              disabled={busy || (isSelf && user.isAdmin)}
              onClick={() => run(() => onUpdate({ userId: user.id, isAdmin: !user.isAdmin }))}
              className={`px-3 py-1.5 rounded-lg text-xs font-medium border cursor-pointer disabled:opacity-30 disabled:cursor-default ${
                user.isAdmin
                  ? "border-amber-500/30 text-amber-300 hover:bg-amber-500/10"
                  : "border-violet-500/30 text-violet-300 hover:bg-violet-500/10"
              }`}
            >
              {user.isAdmin ? "ถอดสิทธิ์" : "ให้สิทธิ์"}
            </button>
          </div>

          <div className="border-t border-white/5 pt-4">
            <div className="text-red-400 text-xs font-medium mb-1">ลบบัญชีถาวร</div>
            <p className="text-xs text-gray-500 mb-2">
              ลบคะแนน แมตช์เดี่ยว และแมตช์ Ranked ทั้งหมดของบัญชีนี้ รวมถึงแมตช์
              Ranked ที่ปรากฏในประวัติของคู่แข่งด้วย กู้คืนไม่ได้ —
              พิมพ์ชื่อผู้เล่นเพื่อยืนยัน
            </p>
            <div className="flex gap-2">
              <input
                value={confirmDelete}
                onChange={(e) => setConfirmDelete(e.target.value)}
                placeholder={user.username}
                className="flex-1 bg-[#0a0a1a] border border-white/10 rounded-lg px-3 py-2 text-xs focus:outline-none focus:border-red-500/40"
              />
              <button
                disabled={busy || isSelf || confirmDelete !== user.username}
                onClick={() => run(() => onDelete(user.id))}
                className="px-3 py-2 rounded-lg text-xs font-medium border border-red-500/30 text-red-400 hover:bg-red-500/10 cursor-pointer disabled:opacity-30 disabled:cursor-default"
              >
                ลบ
              </button>
            </div>
          </div>

          {error && <p className="text-xs text-red-400">{error}</p>}
        </div>

        <footer className="px-5 py-4 border-t border-white/5 flex justify-end gap-2">
          <button
            onClick={onClose}
            className="px-4 py-2 rounded-lg text-sm text-gray-400 hover:text-white cursor-pointer"
          >
            ยกเลิก
          </button>
          <button
            disabled={busy}
            onClick={save}
            className="px-4 py-2 rounded-lg text-sm font-medium bg-violet-500/20 border border-violet-500/30 text-violet-300 hover:bg-violet-500/30 cursor-pointer disabled:opacity-40"
          >
            บันทึก
          </button>
        </footer>
      </div>
    </div>
  );
}
