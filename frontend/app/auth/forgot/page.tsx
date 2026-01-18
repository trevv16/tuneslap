import { Metadata } from "next";
import ForgotClient from "./ForgotClient";

export const metadata: Metadata = {
  title: "Forgot Password",
  description: "Reset your TuneSlap password",
};

export default function ForgotPasswordPage() {
  return <ForgotClient />;
}
