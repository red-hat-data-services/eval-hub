#!/bin/bash
#
# CVE Fix Script for eval-hub
# Generated: 2026-04-01
# Repository: https://github.com/red-hat-data-services/eval-hub
#
# This script applies fixes for all detected CVE vulnerabilities
#

set -e

echo "=========================================="
echo "eval-hub CVE Remediation Script"
echo "=========================================="
echo ""

# Check if we're in the right directory
if [ ! -f "package.json" ]; then
    echo "Error: package.json not found. Please run this script from the eval-hub root directory."
    exit 1
fi

echo "✓ Found package.json"
echo ""

# Backup current state
echo "Creating backup of package-lock.json..."
cp package-lock.json package-lock.json.backup.$(date +%Y%m%d_%H%M%S)
echo "✓ Backup created"
echo ""

# Show current vulnerabilities
echo "Current vulnerability status:"
echo "----------------------------------------"
npm audit --json | jq -r '.metadata.vulnerabilities | "High: \(.high), Moderate: \(.moderate), Low: \(.low), Critical: \(.critical)"'
echo ""

# Step 1: Run automated fix
echo "Step 1: Running automated npm audit fix..."
echo "----------------------------------------"
if npm audit fix; then
    echo "✓ Automated fixes applied successfully"
else
    echo "⚠ Some issues require manual intervention"
fi
echo ""

# Step 2: Check remaining vulnerabilities
echo "Step 2: Checking for remaining vulnerabilities..."
echo "----------------------------------------"
REMAINING=$(npm audit --json | jq -r '.metadata.vulnerabilities.total')
echo "Remaining vulnerabilities: $REMAINING"
echo ""

if [ "$REMAINING" -gt 0 ]; then
    echo "Step 3: Applying manual fixes for breaking changes..."
    echo "----------------------------------------"

    # Fix cucumber-html-reporter (requires downgrade)
    echo "Downgrading cucumber-html-reporter to v6.0.0 (breaking change)..."
    if npm install cucumber-html-reporter@6.0.0; then
        echo "✓ cucumber-html-reporter downgraded successfully"
    else
        echo "✗ Failed to downgrade cucumber-html-reporter"
        echo "  Please run manually: npm install cucumber-html-reporter@6.0.0"
    fi
    echo ""
fi

# Step 4: Final audit
echo "Step 4: Final vulnerability check..."
echo "----------------------------------------"
npm audit
echo ""

# Step 5: Summary
echo "=========================================="
echo "Remediation Complete"
echo "=========================================="
echo ""
echo "Next steps:"
echo "1. Review the changes in package.json and package-lock.json"
echo "2. Run your test suite to ensure nothing broke"
echo "3. Test the cucumber-html-reporter functionality (breaking change)"
echo "4. Commit the changes:"
echo "   git add package.json package-lock.json"
echo "   git commit -m 'fix: resolve CVE vulnerabilities in npm dependencies'"
echo ""
echo "If you need to restore the backup:"
echo "   mv package-lock.json.backup.* package-lock.json"
echo "   npm install"
echo ""
