import api from "./api";
import TokenService from "./token.service";

class AuthService {
    login(user) {
        return api.post('auth/login', {
            username: user.username,
            password: user.password
        })
            .then(response => {
                if (response.data.token) {
                    TokenService.setUser(response.data);
                }
                return response.data;
            })
    }

    logout() {
        TokenService.removeUser();
    }
}

export default new AuthService();
