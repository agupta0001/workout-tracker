import type { AxiosInstance } from "axios";
import axios, { AxiosError } from "axios";
import { firebaseAuth } from "../utils/firebase";

class Http {
  axios: AxiosInstance;

  constructor() {
    this.axios = axios.create({
      baseURL: import.meta.env.VITE_API_BASE_URL,
    });
  }

  protected getToken(): Promise<string> {
    const user = firebaseAuth.currentUser;

    if (!user) throw new Error("User not found");

    return user.getIdToken();
  }

  protected async getTokenHeader(useFirebaseToken: boolean) {
    let token: string;
    if (useFirebaseToken) {
      token = await this.getToken();
      token = `Bearer ${token}`;
    } else {
      token = `Bearer ${localStorage.getItem("jwtToken")}`;
    }
    return {
      Authorization: token,
    };
  }

  #handleError(error: unknown): never {
    if (error instanceof AxiosError && error.response) {
      throw error.response.data;
    }
    throw new Error(
      error instanceof Error ? error.message : "An error occurred"
    );
  }

  async get<T>(
    url: string,
    queryParams: { params?: object } = {},
    useFirebaseToken: boolean = false,
    options: { signal?: AbortSignal } = {}
  ): Promise<T> {
    const tokenHeader = await this.getTokenHeader(useFirebaseToken);
    try {
      const response = await this.axios.get(url, {
        ...queryParams,
        headers: {
          ...tokenHeader,
        },
        ...options,
      });
      return response.data as T;
    } catch (error: unknown) {
      this.#handleError(error);
    }
  }

  async post<T>(
    url: string,
    data: object = {},
    useFirebaseToken: boolean = false
  ): Promise<T> {
    const tokenHeader = await this.getTokenHeader(useFirebaseToken);
    try {
      return await this.axios.post(
        url,
        {
          ...data,
        },
        {
          headers: {
            ...tokenHeader,
          },
        }
      );
    } catch (error: unknown) {
      this.#handleError(error);
    }
  }

  async put<T>(
    url: string,
    data: object = {},
    useFirebaseToken: boolean = false
  ): Promise<T> {
    const tokenHeader = await this.getTokenHeader(useFirebaseToken);
    try {
      return await this.axios.put(
        url,
        {
          ...data,
        },
        {
          headers: {
            ...tokenHeader,
          },
        }
      );
    } catch (error: unknown) {
      this.#handleError(error);
    }
  }

  async patch<T>(
    url: string,
    data: object = {},
    useFirebaseToken: boolean = false
  ): Promise<T> {
    const tokenHeader = await this.getTokenHeader(useFirebaseToken);
    try {
      return await this.axios.patch(
        url,
        {
          ...data,
        },
        {
          headers: {
            ...tokenHeader,
          },
        }
      );
    } catch (error: unknown) {
      this.#handleError(error);
    }
  }

  async delete<T>(
    url: string,
    data?: object,
    useFirebaseToken: boolean = false
  ): Promise<T> {
    const tokenHeader = await this.getTokenHeader(useFirebaseToken);
    try {
      return await this.axios.delete(url, {
        headers: {
          ...tokenHeader,
        },
        data,
      });
    } catch (error: unknown) {
      this.#handleError(error);
    }
  }
}

export default new Http();
