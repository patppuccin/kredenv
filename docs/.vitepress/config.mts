import { defineConfig } from "vitepress";
import type { HeadConfig, DefaultTheme } from "vitepress";

const docsSidebarConfig: DefaultTheme.SidebarItem[] = [
  {
    text: "Getting Started",
    items: [
      { text: "Introduction", link: "/docs/" },
      { text: "Installation", link: "/docs/installation" },
      { text: "Quick Start", link: "/docs/quick-start" },
    ],
  },
  {
    text: "Core Concepts",
    items: [
      { text: "How it Works", link: "/docs/how-it-works" },
      { text: "The Vault", link: "/docs/vault" },
      { text: "Kredsfile Manifest", link: "/docs/kredsfile" },
      { text: "Namespaces", link: "/docs/namespaces" },
    ],
  },
  {
    text: "Commands",
    items: [
      { text: "Setup Commands", link: "/docs/setup-commands" },
      { text: "Environment Commands", link: "/docs/environment-commands" },
      { text: "Secrets Commands", link: "/docs/secrets-commands" },
    ],
  },
  {
    text: "Integrations",
    items: [
      { text: "Shell Hooks", link: "/docs/shell-hooks" },
      { text: "Prompt Frameworks", link: "/docs/prompt-frameworks" },
      { text: "External Stores", link: "/docs/external-stores" },
    ],
  },
  {
    text: "Reference",
    items: [{ text: "Caveats", link: "/docs/caveats" }],
  },
];

const socialLinksConfig: DefaultTheme.SocialLink[] = [
  {
    icon: "github",
    link: "https://github.com/patppuccin/kredenv",
    ariaLabel: "GitHub",
  },
];

const navBarConfig: DefaultTheme.NavItem[] = [
  { text: "Docs", link: "/docs/" },
  { text: "About", link: "/about" },
  {
    text: "Extras",
    items: [
      { text: "License", link: "/license" },
      { text: "Roadmap", link: "/roadmap" },
    ],
  },
];

const htmlAttrs: HeadConfig[] = [
  [
    "link",
    {
      rel: "icon",
      type: "image/png",
      href: "/favicon-96x96.png",
      sizes: "96x96",
    },
  ],
  ["link", { rel: "icon", type: "image/svg+xml", href: "/favicon.svg" }],
  ["link", { rel: "shortcut icon", href: "/favicon.ico" }],
  [
    "link",
    {
      rel: "apple-touch-icon",
      sizes: "180x180",
      href: "/apple-touch-icon.png",
    },
  ],
  ["meta", { name: "apple-mobile-web-app-title", content: "Kredenv" }],
  ["link", { rel: "manifest", href: "/site.webmanifest" }],
];

// https://vitepress.dev/reference/site-config
export default defineConfig({
  srcDir: "content",
  head: htmlAttrs,

  title: "kredenv",
  description: "Inject env vars & secrets into your shell environment",
  themeConfig: {
    logo: "/favicon.svg",
    siteTitle: "Kredenv",
    search: { provider: "local" },
    outline: [2, 3],
    nav: navBarConfig,
    sidebar: { "/docs/": docsSidebarConfig },
    socialLinks: socialLinksConfig,
    footer: {
      message: "Released under the <a href='/license'>Apache 2.0 License</a>",
      copyright: `Copyright © ${new Date().getFullYear()} <a href='https://patrickambrose.com'>Patrick Ambrose</a>.`,
    },
  },
  markdown: {
    theme: {
      light: "vitesse-light",
      dark: "vitesse-dark",
    },
  },
});
