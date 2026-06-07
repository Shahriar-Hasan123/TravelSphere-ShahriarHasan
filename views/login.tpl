<div class="page-body">
    <div class="login-wrap">
        <h1 class="page-title">Sign in</h1>
        {{if .Error}}
        <div class="alert alert-error">{{.Error}}</div>
        {{end}}
        <form class="login-form" method="POST" action="/login">
            <div class="form-group">
                <label class="form-label">Username</label>
                <input type="text" name="username" class="form-input"
                       required autocomplete="username">
            </div>
            <div class="form-group">
                <label class="form-label">Password</label>
                <input type="password" name="password" class="form-input"
                       required autocomplete="current-password">
            </div>
            <button type="submit" class="btn-primary">Login</button>
        </form>
    </div>
</div>