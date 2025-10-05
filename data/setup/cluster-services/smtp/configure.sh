#!/bin/bash

print_header "Setting up SMTP Configuration"

print_info "SMTP configuration allows Wild Cloud applications to send transactional emails"
print_info "(password resets, notifications, etc.) through your email service provider."
echo ""

# Collect SMTP configuration
print_info "Collecting SMTP configuration..."
prompt_if_unset_config "cloud.smtp.host" "Enter SMTP host (e.g., email-smtp.us-east-2.amazonaws.com for AWS SES)" ""
prompt_if_unset_config "cloud.smtp.port" "Enter SMTP port (usually 465 for SSL, 587 for STARTTLS)" "465"
prompt_if_unset_config "cloud.smtp.user" "Enter SMTP username/access key" ""
prompt_if_unset_config "cloud.smtp.from" "Enter default 'from' email address" "no-reply@$(wild-config cloud.domain)"
prompt_if_unset_config "cloud.smtp.tls" "Enable TLS? (true/false)" "true"
prompt_if_unset_config "cloud.smtp.startTls" "Enable STARTTLS? (true/false)" "true"
print_success "SMTP configuration collected successfully"

# Collect SMTP password/secret
print_info "Setting up SMTP password..."
echo ""
echo "For AWS SES, this would be your Secret Access Key."
echo "For Gmail/Google Workspace, this would be an App Password."
echo "For other providers, this would be your SMTP password."
echo ""

prompt_if_unset_secret "cloud.smtp.password" "Enter SMTP password/secret key" ""

print_success "SMTP configuration setup complete!"
echo ""
echo "Your SMTP settings:"
echo "  Host: $(wild-config cloud.smtp.host)"
echo "  Port: $(wild-config cloud.smtp.port)"
echo "  User: $(wild-config cloud.smtp.user)"
echo "  From: $(wild-config cloud.smtp.from)"
echo "  Password: $(wild-secret cloud.smtp.password >/dev/null 2>&1 && echo "✓ Set" || echo "✗ Not set")"
echo ""
