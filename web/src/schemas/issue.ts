import { z } from "zod";

export const StatusType = z.enum(["Incomplete", "In-Progress", "Complete"]);
export type StatusType = z.infer<typeof StatusType>;

export const IssueSchema = z.object({
  id: z.number(),
  user_id: z.uuid(),
  user_name: z.string(),
  title: z.string().min(1, "Title is required").max(255, "Title too long"),
  description: z.string(),
  status: StatusType,
  created_at: z.string(),
  updated_at: z.string(),
});
export type Issue = z.infer<typeof IssueSchema>;

export const CreateIssueSchema = z.object({
  title: z.string().min(1, "Title is required").max(255, "Title too long"),
  description: z.string().default(""),
});
export type CreateIssueInput = z.infer<typeof CreateIssueSchema>;

export const UpdateStatusSchema = z.object({
  status: StatusType,
});
export type UpdateStatusInput = z.infer<typeof UpdateStatusSchema>;
