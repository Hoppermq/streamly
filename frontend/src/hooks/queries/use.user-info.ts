import {useAuthStore} from "@/hooks/store/auth.store.ts";
import {useQuery} from "@tanstack/react-query";

const ZITADEL_ISSUER = ""

export const useUserInfo = () => {
  const user = useAuthStore((state) => state.user);

  return useQuery({
    queryKey: ['user-info', user?.profile.sub],
    queryFn: async () => {
      const response = await fetch(`${ZITADEL_ISSUER}/oidc/v1/userinfo`, {
        headers: {
          Authorization: `Bearer ${user?.access_token}`,
        }
      })
      return response.json();
    },
    enabled: !!user,
    staleTime: 5 * 60 * 1000,
  })
}
