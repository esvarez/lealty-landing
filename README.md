# Astro Starter Kit: Basics

```sh
npm create astro@latest -- --template basics
```

[![Open in StackBlitz](https://developer.stackblitz.com/img/open_in_stackblitz.svg)](https://stackblitz.com/github/withastro/astro/tree/latest/examples/basics)
[![Open with CodeSandbox](https://assets.codesandbox.io/github/button-edit-lime.svg)](https://codesandbox.io/p/sandbox/github/withastro/astro/tree/latest/examples/basics)
[![Open in GitHub Codespaces](https://github.com/codespaces/badge.svg)](https://codespaces.new/withastro/astro?devcontainer_path=.devcontainer/basics/devcontainer.json)

> 🧑‍🚀 **Seasoned astronaut?** Delete this file. Have fun!

![just-the-basics](https://github.com/withastro/astro/assets/2244813/a0a5533c-a856-4198-8470-2d67b1d7c554)

## Mailchimp Integration Setup

To enable email collection with Mailchimp:

1. **Get your Mailchimp API Key:**
   - Log in to your Mailchimp account
   - Go to Account → Extras → API keys
   - Create a new API key or copy an existing one

2. **Find your Server Prefix:**
   - Look at the end of your API key (e.g., `abc123-us1`)
   - The server prefix is the part after the dash (e.g., `us1`)

3. **Get your Audience ID:**
   - Go to Audience → All contacts → Settings → Audience name and defaults
   - Copy the Audience ID

4. **Set up environment variables:**
   - Copy `.env.example` to `.env`
   - Fill in your Mailchimp credentials:
     ```
     MAILCHIMP_API_KEY=your_api_key_here
     MAILCHIMP_SERVER_PREFIX=us1
     MAILCHIMP_AUDIENCE_ID=your_audience_id_here
     ```

5. **Test the integration:**
   - Start the development server: `npm run dev`
   - Try subscribing with a test email address
   - Check your Mailchimp audience to confirm the subscriber was added

## 🚀 Project Structure

Inside of your Astro project, you'll see the following folders and files:

```text
/
├── public/
│   └── favicon.svg
├── src/
│   ├── layouts/
│   │   └── Layout.astro
│   └── pages/
│       └── index.astro
└── package.json
```

To learn more about the folder structure of an Astro project, refer to [our guide on project structure](https://docs.astro.build/en/basics/project-structure/).

## 🧞 Commands

All commands are run from the root of the project, from a terminal:

| Command                   | Action                                           |
| :------------------------ | :----------------------------------------------- |
| `npm install`             | Installs dependencies                            |
| `npm run dev`             | Starts local dev server at `localhost:4321`      |
| `npm run build`           | Build your production site to `./dist/`          |
| `npm run preview`         | Preview your build locally, before deploying     |
| `npm run astro ...`       | Run CLI commands like `astro add`, `astro check` |
| `npm run astro -- --help` | Get help using the Astro CLI                     |

## 👀 Want to learn more?

Feel free to check [our documentation](https://docs.astro.build) or jump into our [Discord server](https://astro.build/chat).
