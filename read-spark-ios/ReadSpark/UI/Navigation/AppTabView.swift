import SwiftUI

struct AppTabView: View {
    var body: some View {
        TabView {
            HomeView()
                .tabItem {
                    Label("首页", systemImage: "house")
                }

            CategoryView()
                .tabItem {
                    Label("分类", systemImage: "square.grid.2x2")
                }

            Text("阅读列表")
                .tabItem {
                    Label("阅读", systemImage: "book")
                }

            DiscoverView()
                .tabItem {
                    Label("发现", systemImage: "compass")
                }

            ProfileView()
                .tabItem {
                    Label("我的", systemImage: "person")
                }
        }
    }
}
