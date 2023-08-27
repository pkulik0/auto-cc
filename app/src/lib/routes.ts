interface Route {
    label: string,
    destination: string
}

export const routes: Route[] = [
    { label: "Videos", destination: "/" },
    { label: "Settings", destination: "/settings" }
]