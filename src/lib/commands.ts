export type CommandDef = {
  command: string;
  description: string;
};

export const COMMANDS: CommandDef[] = [
  { command: "/login", description: "Save an API key profile" },
  { command: "/logout", description: "Clear the active profile" },
  { command: "/whoami", description: "Show organisation for the current key" },
  { command: "/org get", description: "Organisation profile" },
  { command: "/org billing", description: "Billing entitlements" },
  { command: "/org stats", description: "Organisation stats" },
  { command: "/assessments list", description: "List assessments" },
  { command: "/assessments get", description: "Get assessment by slug" },
  { command: "/invites list", description: "List invitations" },
  { command: "/invites create", description: "Create an invitation" },
  { command: "/results list", description: "List results for an assessment" },
  { command: "/tasks list", description: "List org tasks" },
  { command: "/pipelines list", description: "List pipelines" },
  { command: "/webhooks list", description: "List webhooks" },
  { command: "/interviews list", description: "List interviews" },
  { command: "/integrations list", description: "List ATS integrations" },
  { command: "/ai", description: "Ask the Assess assistant" },
  { command: "/help", description: "Show help" },
  { command: "/clear", description: "Clear the message feed" },
  { command: "/exit", description: "Quit the CLI" },
];
