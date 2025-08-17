# Email Debugging Guide

## ✅ Progress Update

- **TLS Connection**: ✅ Working perfectly
- **SMTP Server**: ✅ smtp.zoho.com:587 is accessible
- **Issue**: ❌ Authentication Failed (535 error)

## Quick Test

Run the email test to see detailed logs:

```bash
cd debug
go run test_email.go
```

## 🔍 Current Issue: Authentication Failed (535)

### Zoho SMTP Authentication Requirements

1. **Enable IMAP Access**

   - Log into your Zoho Mail account
   - Go to Settings → Mail → POP/IMAP
   - Enable IMAP access

2. **Use App Password (Recommended)**

   - Go to Zoho Account Settings → Security → App Passwords
   - Generate a new App Password for "Mail"
   - Use this App Password instead of your regular password in the `.env` file

3. **Enable Less Secure Apps (Alternative)**

   - In Zoho Account Settings → Security
   - Enable "Allow less secure apps" (not recommended)

4. **Two-Factor Authentication**
   - If 2FA is enabled, you MUST use App Passwords
   - Regular passwords won't work with 2FA enabled

### Quick Fixes to Try

1. **Verify Credentials**

   ```bash
   # Make sure these are correct in your .env file:
   MAIL_USERNAME=noreply@songbooksofpraise.com  # Must be full email
   MAIL_PASSWORD=your_app_password_here         # Use App Password
   ```

2. **Try Manual Login**

   - Try logging into Zoho Mail with the same credentials
   - If that fails, the credentials are wrong

3. **Check Account Status**
   - Make sure the account is active
   - Verify the domain is properly configured in Zoho

## Debugging Steps

### 1. Check Environment Variables

Make sure these are set in your `.env` file:

- `MAIL_HOST=smtp.zoho.com`
- `MAIL_PORT=587`
- `MAIL_USERNAME=noreply@songbooksofpraise.com`
- `MAIL_PASSWORD=your_password`

### 2. Common Issues

#### Authentication Errors (535)

- ✅ **Most Likely**: Need to use App Password instead of regular password
- ✅ **Check**: IMAP access is enabled in Zoho settings
- ✅ **Verify**: Username is the full email address
- ✅ **Confirm**: Account has SMTP access enabled

#### Connection Errors (Fixed ✅)

- ✅ Port 587 (TLS) is now working correctly
- ✅ STARTTLS is properly implemented

### 3. Test with Different Ports

If port 587 doesn't work, try port 465 with SSL:

```bash
# In .env file
MAIL_PORT=465
```

### 4. Manual SMTP Test

Test authentication manually:

```bash
# Test connection and auth
telnet smtp.zoho.com 587
# Then follow SMTP commands to test auth
```

## Common Error Messages

- ✅ **"failed to connect to SMTP server"** - Fixed!
- ❌ **"535 Authentication Failed"** - Current issue (see solutions above)
- **"failed to set sender"** - Email address not authorized
- **"failed to set recipient"** - Invalid recipient email

## Next Steps

1. ✅ **Connection Working** - STARTTLS is functioning correctly
2. ❌ **Fix Authentication** - Get App Password from Zoho
3. **Update .env** - Use App Password instead of regular password
4. **Test Again** - Run the test program
5. **Check spam folders** - Once authentication works

## 🎯 Immediate Action Required

**Generate Zoho App Password:**

1. Go to https://accounts.zoho.com/home#security/apppassword
2. Click "Generate New Password"
3. Select "Mail" as the application
4. Copy the generated password
5. Update your `.env` file with this App Password

This should resolve the authentication issue!
