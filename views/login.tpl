<div class="page-body">
    <div class="login-wrap">
        <h1 class="page-title">Sign in</h1>
        <p class="page-subtitle" style="margin-bottom:1.5rem;">
            Use username <strong>beta</strong> and password <strong>beta123</strong>.
        </p>

        {{if .Error}}
        <div class="alert alert-error">{{.Error}}</div>
        {{end}}

        <form class="login-form" method="POST" action="/login">
            <!-- Carry the redirect target through the form submission -->
            <input type="hidden" name="redirect_to" value="{{.RedirectTo}}">

            <div class="form-group">
                <label class="form-label" for="username">Username</label>
                <input
                    type="text"
                    id="username"
                    name="username"
                    class="form-input"
                    value="{{.Username}}"
                    required
                    autocomplete="username"
                    autofocus
                >
            </div>

            <div class="form-group">
                <label class="form-label" for="password">Password</label>
                <input
                    type="password"
                    id="password"
                    name="password"
                    class="form-input"
                    required
                    autocomplete="current-password"
                >
            </div>

            <button type="submit" class="btn-primary" style="width:100%;">
                Login
            </button>
        </form>
    </div>
</div>