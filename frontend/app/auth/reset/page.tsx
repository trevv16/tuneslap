import { Metadata } from "next";
import ResetClient from "./ResetClient";

export const metadata: Metadata = {
  title: "Reset Password",
  description: "Reset your TuneSlap password",
};

export default function ResetPasswordPage() {
  return <ResetClient />;
}
