export function useUser(): [string, Error | null] {
    return [localStorage.getItem('email') ?? '', null]
}