# Lessons Learned - SupaManager Development

## Session: API Response Format Debugging (Nov 16, 2025)

### Problem Summary
The Supabase Studio editor page was showing "Cannot read properties of undefined (reading 'data')" error, while the main project page loaded correctly.

### Root Cause
The `postPlatformPgMetaQuery` API endpoint was returning data in the wrong format. The `executeSql` function in Studio automatically wraps API responses in `{ result: data }`, but we were pre-wrapping the response, causing double-wrapping.

### The Investigation Journey

#### Mistake #1: Assuming the Issue Was Missing Fields
**What happened:** Initially thought Studio required specific fields (`dbVersion`, `connectionString`, `restUrl`, `kpsVersion`, `url`) that were missing from the API response.

**What went wrong:** Added these fields with values like `https://project-ref.supamanager.io`, but the project ref contained an apostrophe (`volatilize-trolleybus's`), creating invalid URLs. This caused different errors.

**What we learned:**
- Always test with the **actual data** (including special characters like apostrophes)
- Don't assume missing fields - verify what Studio actually needs by checking the source code
- Invalid hostnames (containing apostrophes) cause URL parsing to fail

#### Mistake #2: Not Understanding the Architecture
**What happened:** Made changes trying to fix errors without understanding how Studio's main page worked without those fields.

**What went wrong:** The main page worked fine on `origin/main` branch (which had empty/mock data), but broke when we added the "required" fields.

**What we learned:**
- **Always compare with a working baseline** - the origin/main code was working, which proved our additions broke it
- Studio has conditional logic - it doesn't require certain fields if they're not present
- User's suggestion to "compare with main branch" was crucial and saved time

#### Mistake #3: Not Understanding Response Wrapping
**What happened:** The `postPlatformPgMetaQuery` endpoint returned:
```json
{"result": [{"data": {"entities": [], "count": 0}}]}
```

**What went wrong:** The `executeSql` function in Studio does:
```typescript
const { data, error } = await post('/platform/pg-meta/{ref}/query', ...)
return { result: data }
```

So our pre-wrapped response became:
```json
{
  "result": {
    "result": [{"data": {"entities": [], "count": 0}}]
  }
}
```

Then `result[0]` tried to access the first element of an object (not array), returning `undefined`.

**What we learned:**
- **Understand the full data flow** - don't just look at one endpoint
- Check how the client-side code processes responses
- Look at **other working endpoints** to see the expected format (like `authorization`, `repositories`)
- API response wrapping/transformation happens at multiple layers

### Correct Solution
Return the data array **directly** from the API:
```json
[{"data": {"entities": [], "count": 0}}]
```

Let Studio's `executeSql` wrapper handle the `{ result: data }` wrapping.

### Key Debugging Techniques That Worked

1. **Check the Network Tab**
   - See exactly what the API is returning
   - Compare with other working endpoints
   - Inspect actual HTTP requests/responses

2. **Compare with Working Code**
   - Test `origin/main` branch to see if it works
   - Use git diff to understand what changed

3. **Search the Source Code**
   - Use grep to find where errors occur: `grep -rn "split()" apps/studio/`
   - Find the actual code that processes the API response

4. **Check API Logs**
   - Use `docker compose logs` to see what requests are being made
   - Verify the endpoint is actually being called

5. **Test Directly with curl**
   - Bypass the UI to test API responses directly
   - Compare actual response with expected format

### Mistakes to NEVER Repeat

#### 1. Don't Guess at Required Fields
**Bad:** "Studio probably needs these fields, let me add them"
**Good:** "Let me check Studio's source code to see exactly what it expects"

**How to do it right:**
```bash
# Find where the error occurs
grep -rn "Cannot read properties" apps/studio/

# Find what fields are being accessed
grep -rn "lastPage.data" apps/studio/data/
```

#### 2. Don't Assume Response Formats
**Bad:** "APIs should return `{ result: data }`"
**Good:** "Let me check what format other endpoints use"

**How to do it right:**
- Look at working endpoints (authorization, repositories, etc.)
- Check the client-side fetcher code to see what it expects
- Test the exact response format with curl

#### 3. Don't Ignore Special Characters in Test Data
**Bad:** Generating URLs with unvalidated user input (`project-ref's.domain.com`)
**Good:** Validate/sanitize data before using in URLs, or check if fields are actually needed

**How to do it right:**
```go
// If returning URLs, sanitize the input
safeRef := url.PathEscape(projectRef)
url := fmt.Sprintf("https://%s.domain.com", safeRef)

// OR better: don't return fields that aren't needed
// Check if Studio actually requires them first
```

#### 4. Don't Change Multiple Things at Once
**Bad:** Added 5+ fields trying to fix one error
**Good:** Add one field, test, add another, test

**How to do it right:**
- Make minimal changes
- Test after each change
- Use git to track what you changed
- Commit small, focused changes

#### 5. Always Check the Working Baseline First
**Bad:** "It's broken, let me start fixing things"
**Good:** "Let me check if `origin/main` works, then compare"

**How to do it right:**
```bash
# Save your changes
git stash

# Test the original code
git checkout origin/main
# Test in browser

# If it works, your changes broke it
# Compare to find what changed
git diff origin/main
```

### Architecture Insights

#### How Studio's `executeSql` Works
```typescript
// 1. Studio calls executeSql
const { result } = await executeSql({
  projectRef,
  sql: "SELECT ..."
})

// 2. executeSql makes HTTP request
const { data, error } = await post('/platform/pg-meta/{ref}/query', {
  body: { query: sql }
})

// 3. executeSql wraps the response
return { result: data }  // ← Wraps whatever API returns

// 4. Caller accesses result[0]
const row = result[0]  // ← Expects result to be an array
```

**Key insight:** The API should return an **array**, because:
- `executeSql` wraps it: `{ result: [array] }`
- Caller does: `result[0]`
- Must be an array for `[0]` to work

#### How to Debug API Response Formats
1. **Find the client-side call:**
   ```bash
   grep -rn "executeSql" apps/studio/data/
   ```

2. **Check what it does with the result:**
   ```typescript
   const { result } = await executeSql(...)
   return result[0]  // ← Expects array
   ```

3. **Check the fetcher/wrapper:**
   ```typescript
   // In execute-sql-query.ts
   const { data } = await post(...)
   return { result: data }  // ← Wraps response
   ```

4. **Determine correct API format:**
   - If wrapper does `{ result: data }` and caller does `result[0]`
   - API should return: `[{...}]` (array)
   - NOT: `{ result: [{...}] }` (pre-wrapped)

### Documentation Updates Needed

- [x] Created LESSONS_LEARNED.md with detailed debugging notes
- [ ] Update ARCHITECTURE.md with API response format guidelines
- [ ] Add troubleshooting guide for common Studio integration issues
- [ ] Document the pg-meta query endpoint format

### Future Improvements

1. **Add API Response Validation**
   - Create tests that verify response formats match Studio expectations
   - Add runtime validation/logging for debugging

2. **Better Error Messages**
   - Return structured errors from API
   - Log actual vs expected formats when mismatches occur

3. **Integration Tests**
   - Test Studio UI against the API
   - Catch format mismatches early

4. **Code Documentation**
   - Comment complex response wrapping logic
   - Document expected formats in API handlers

### Summary: What Actually Fixed It
✅ Return raw array from API: `[{"data": {...}}]`
✅ Let executeSql handle wrapping: `{ result: [{"data": {...}}] }`
✅ Caller gets correct format: `result[0]` = `{"data": {...}}`

### Time Spent
- Investigating: 2-3 hours
- Failed attempts with field additions: 1 hour
- Finding root cause (response wrapping): 30 minutes
- Implementing fix: 10 minutes
- **Total: ~4 hours**

**Lesson:** Could have been 30 minutes if we:
1. Checked working baseline first
2. Looked at other endpoint formats
3. Understood the response wrapping flow

---

*Generated on: 2025-11-16*
*Session: API debugging and fix*
