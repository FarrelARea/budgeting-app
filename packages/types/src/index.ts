import { z } from "zod";

// ---------- primitives ----------
export const uuid = z.string().uuid();
export const isoDate = z.string().datetime({ offset: true });
export const money = z.number().finite(); // amounts are in minor units (e.g. cents) on the wire

// ---------- domain ----------
export const AccountSchema = z.object({
  id: uuid,
  name: z.string().min(1).max(80),
  currency: z.string().length(3),
  balance: money,
  createdAt: isoDate,
  updatedAt: isoDate
});
export type Account = z.infer<typeof AccountSchema>;

export const CategoryKindSchema = z.enum(["income", "expense"]);
export type CategoryKind = z.infer<typeof CategoryKindSchema>;

export const CategorySchema = z.object({
  id: uuid,
  name: z.string().min(1).max(60),
  kind: CategoryKindSchema,
  color: z.string().regex(/^#[0-9a-fA-F]{6}$/).optional(),
  parentId: uuid.nullable().optional()
});
export type Category = z.infer<typeof CategorySchema>;

export const TransactionSchema = z.object({
  id: uuid,
  accountId: uuid,
  categoryId: uuid.nullable(),
  amount: money, // signed: positive = income, negative = expense
  description: z.string().max(240).default(""),
  occurredAt: isoDate,
  createdAt: isoDate,
  updatedAt: isoDate
});
export type Transaction = z.infer<typeof TransactionSchema>;

export const BudgetPeriodSchema = z.enum(["weekly", "monthly", "yearly"]);
export type BudgetPeriod = z.infer<typeof BudgetPeriodSchema>;

export const BudgetSchema = z.object({
  id: uuid,
  categoryId: uuid,
  amount: money,
  period: BudgetPeriodSchema,
  startsAt: isoDate,
  endsAt: isoDate.nullable().optional()
});
export type Budget = z.infer<typeof BudgetSchema>;

// ---------- API envelope ----------
export const ApiErrorSchema = z.object({
  code: z.string(),
  message: z.string(),
  details: z.unknown().optional()
});
export type ApiError = z.infer<typeof ApiErrorSchema>;

export const ApiResponseSchema = <T extends z.ZodTypeAny>(data: T) =>
  z.discriminatedUnion("ok", [
    z.object({ ok: z.literal(true), data }),
    z.object({ ok: z.literal(false), error: ApiErrorSchema })
  ]);

export type ApiResponse<T> =
  | { ok: true; data: T }
  | { ok: false; error: ApiError };

// ---------- request payloads ----------
export const CreateTransactionInputSchema = TransactionSchema.pick({
  accountId: true,
  categoryId: true,
  amount: true,
  description: true,
  occurredAt: true
});
export type CreateTransactionInput = z.infer<typeof CreateTransactionInputSchema>;

export const CreateAccountInputSchema = AccountSchema.pick({
  name: true,
  currency: true
});
export type CreateAccountInput = z.infer<typeof CreateAccountInputSchema>;
