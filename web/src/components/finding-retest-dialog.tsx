"use client";

import * as React from "react";

import { RotateCcwIcon } from "lucide-react";
import { toast } from "sonner";

import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Field, FieldDescription, FieldGroup, FieldLabel } from "@/components/ui/field";
import { Spinner } from "@/components/ui/spinner";
import { Textarea } from "@/components/ui/textarea";
import { api } from "@/lib/api";
import type { FindingRetest } from "@/lib/types";

interface FindingRetestDialogProps {
  findingId: string;
  findingName?: string;
  onClose: () => void;
  onStarted?: (retest: FindingRetest) => void;
}

// 열 때만 마운트하고 닫은 후 설명을 비웁니다; 목록과 상세에서 제출 잠금 및 오류 처리를 공유하며, 시작 후 현재 페이지에 머무릅니다.
export function FindingRetestDialog({ findingId, findingName, onClose, onStarted }: FindingRetestDialogProps) {
  const notesId = React.useId();
  const [notes, setNotes] = React.useState("");
  const [submitting, setSubmitting] = React.useState(false);
  const submitLock = React.useRef(false);

  async function start() {
    if (submitLock.current) return;
    submitLock.current = true;
    setSubmitting(true);
    try {
      const result = await api.startFindingRetest(findingId, notes.trim());
      onStarted?.(result.retest);
      onClose();
      toast.success(result.created ? "재검증을 시작했습니다. '재검증 중'을 클릭하여 세션을 볼 수 있습니다" : "해당 취약점은 재검증 중입니다. 기존 세션을 볼 수 있습니다");
    } catch (e) {
      toast.error(`재검증 시작 실패: ${(e as Error).message}`);
    } finally {
      submitLock.current = false;
      setSubmitting(false);
    }
  }

  return (
    <Dialog open onOpenChange={(open) => !open && !submitLock.current && onClose()}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>취약점 재검증 #{findingId}</DialogTitle>
          <DialogDescription className="break-words">
            {findingName ? <span className="mb-2 block">{findingName}</span> : null}
            재검증 Agent가
            기존 증거와 테스트 제약 조건을 읽고 독립된 세션에서 해당 취약점을 검증합니다. 재검증이 성공적으로 완료되고 수정이 확인되면 취약점 상태가 자동으로 '수정 완료'로 변경되며, 그 외 결론에서는 기존 상태를 유지합니다.
          </DialogDescription>
        </DialogHeader>
        <FieldGroup>
          <Field data-disabled={submitting}>
            <FieldLabel htmlFor={notesId}>추가 설명(선택 사항)</FieldLabel>
            <Textarea
              id={notesId}
              value={notes}
              maxLength={4000}
              rows={4}
              disabled={submitting}
              onChange={(e) => setNotes(e.target.value)}
              placeholder="예: 기존 테스트 계정으로 기존 엔드포인트를 검증하세요; 수정 버전은 v2입니다."
            />
            <FieldDescription>수정 버전·테스트 조건 또는 이번 재검증의 제한 사항을 추가할 수 있습니다.</FieldDescription>
          </Field>
        </FieldGroup>
        <DialogFooter>
          <Button variant="outline" disabled={submitting} onClick={onClose}>
            취소
          </Button>
          <Button disabled={submitting} onClick={() => void start()}>
            {submitting ? <Spinner data-icon="inline-start" /> : <RotateCcwIcon data-icon="inline-start" />}
            {submitting ? "생성 중…" : "재검증 시작"}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
