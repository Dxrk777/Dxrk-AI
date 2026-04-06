<!--
╔═══════════════════════════════════════════════════════════════════════════════╗
║                                                                               ║
║     ██████╗ ██╗  ██╗██████╗ ██╗  ██╗     ██████╗ ███╗   ███╗███████╗       ║
║     ██╔══██╗╚██╗██╔╝██╔══██╗██║ ██╔╝    ██╔═══██╗████╗ ████║██╔════╝       ║
║     ██║  ██║ ╚███╔╝ ██████╔╝█████╔╝     ██║   ██║██╔████╔██║█████╗         ║
║     ██║  ██║ ██╔██╗ ██╔══██╗██╔═██╗     ██║   ██║██║╚██╔╝██║██╔══╝         ║
║     ██████╔╝██╔╝ ██╗██║  ██║██║  ██╗    ╚██████╔╝██║ ╚═╝ ██║███████╗       ║
║     ╚═════╝ ╚═╝  ╚═╝╚═╝  ╚═╝╚═╝  ╚═╝     ╚═════╝ ╚═╝     ╚═╝╚══════╝       ║
║                                                                               ║
║              DXRK AI INTELLIGENCE SYSTEM — SKILL PACKAGE                      ║
║                                                                               ║
╚═══════════════════════════════════════════════════════════════════════════════╝
-->

# Security Review Skill

## Trigger
"security review", "analyze security", "check for vulnerabilities", "OWASP", "security audit"

## Purpose
Perform semantic security analysis detecting OWASP Top 10 vulnerabilities, injection attacks, auth flaws, and data exposure.

## Analysis Areas
- **Injection**: SQL, NoSQL, Command, XSS, LDAP
- **Authentication**: Broken auth, credential handling, session management
- **Sensitive Data**: Exposure in logs, responses, client-side storage
- **Access Control**: IDOR, privilege escalation, broken authorization
- **Security Misconfiguration**: Default credentials, verbose errors, exposed configs

## Process
1. Scan changed files for security anti-patterns
2. Check authentication/authorization logic
3. Verify input sanitization and validation
4. Check for secrets in code
5. Report findings with severity (CRITICAL/HIGH/MEDIUM/LOW)

## Output Format
```
## Security Findings

### [CRITICAL] <title>
**File:** <path>
**Line:** <number>
**Issue:** <description>
**Fix:** <recommendation>
```
