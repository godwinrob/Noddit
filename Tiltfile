# Tiltfile for Noddit local development
# Run with: tilt up

# Allow configuration via tilt_config.json
config.define_string_list("to-run", args=True)
cfg = config.parse()
groups = {
    'backend': ['postgres', 'backend'],
    'frontend': ['frontend'],
}

# Determine what to run
to_run = cfg.get('to-run', ['backend', 'frontend'])
resources = [r for r in to_run for r in groups.get(r, [r])]

# Read .env file for build args (NEXT_PUBLIC_* must be present at build time)
def read_dotenv(path):
    result = {}
    content = str(read_file(path, ''))
    for line in content.splitlines():
        line = line.strip()
        if not line or line.startswith('#'):
            continue
        if '=' in line:
            key, _, value = line.partition('=')
            result[key.strip()] = value.strip()
    return result

env = read_dotenv('./.env')

# Load docker-compose
docker_compose('./docker-compose.yml')

# PostgreSQL Database
if 'postgres' in resources:
    dc_resource(
        'postgres',
        labels=['database'],
        links=[
            link('postgres://postgres:****@localhost:5432/userdb', 'Database Connection'),
        ],
    )

# Go Backend
if 'backend' in resources:
    dc_resource(
        'backend',
        labels=['backend'],
        resource_deps=['postgres'],
        links=[
            link('http://localhost:8080/api/public/recentposts', 'API - Recent Posts'),
        ],
    )

    # Watch for Go file changes and rebuild
    docker_build(
        'noddit-backend',
        context='./backend-go',
        dockerfile='./backend-go/Dockerfile',
        live_update=[
            sync('./backend-go', '/root'),
            run(
                'cd /root && go build -o main ./cmd/api',
                trigger=['./backend-go/**/*.go']
            ),
            restart_container(),
        ],
    )

# Next.js Frontend
if 'frontend' in resources:
    dc_resource(
        'frontend',
        labels=['frontend'],
        resource_deps=['backend'],
        links=[
            link('http://localhost:8081', 'Open Frontend'),
        ],
    )

    # Watch for Next.js file changes
    docker_build(
        'noddit-frontend',
        context='./noddit-next',
        dockerfile='./noddit-next/Dockerfile',
        build_args={
            'NEXT_PUBLIC_API_URL': env.get('NEXT_PUBLIC_API_URL', 'http://localhost:8080'),
            'NEXT_PUBLIC_CLERK_PUBLISHABLE_KEY': env.get('NEXT_PUBLIC_CLERK_PUBLISHABLE_KEY', ''),
        },
        live_update=[
            sync('./noddit-next', '/app'),
            run(
                'cd /app && npm run build',
                trigger=['./noddit-next/app/**/*', './noddit-next/components/**/*', './noddit-next/lib/**/*']
            ),
            restart_container(),
        ],
    )

    # Auto-open browser when frontend is ready
    local_resource(
        'open-frontend',
        cmd='sleep 3 && python -m webbrowser http://localhost:8081 || start http://localhost:8081 || open http://localhost:8081',
        resource_deps=['frontend'],
        labels=['frontend'],
        auto_init=True,
        trigger_mode=TRIGGER_MODE_MANUAL,
    )

# Local development mode (alternative to Docker for faster iteration)
# Uncomment to use local processes instead of containers

# local_resource(
#     'backend-dev',
#     serve_cmd='cd backend-go && go run cmd/api/main.go',
#     deps=['./backend-go'],
#     labels=['backend'],
#     resource_deps=['postgres'],
# )

# local_resource(
#     'frontend-dev',
#     serve_cmd='cd noddit-next && npm run dev',
#     deps=['./noddit-next'],
#     labels=['frontend'],
#     resource_deps=['backend-dev'],
# )

# Port forwards (useful if running services inside k8s)
# k8s_resource('backend', port_forwards='8080:8080')
# k8s_resource('frontend', port_forwards='8081:8081')

print("""
╔══════════════════════════════════════════════════════════╗
║                    Noddit Development                     ║
╠══════════════════════════════════════════════════════════╣
║  Backend:   http://localhost:8080                         ║
║  Frontend:  http://localhost:8081                         ║
║  Database:  postgres://postgres:****@localhost:5432       ║
╠══════════════════════════════════════════════════════════╣
║  Press SPACE to open the Tilt UI                          ║
║  Press 'q' to quit                                        ║
╚══════════════════════════════════════════════════════════╝
""")
