import SwiftUI

struct AppTabView: View {
    var body: some View {
        TabView {
            Text("Home")
                .tabItem {
                    Label("Home", systemImage: "house")
                }

            Text("Profile")
                .tabItem {
                    Label("Profile", systemImage: "person")
                }
        }
    }
}
